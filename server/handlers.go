package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"

	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/Broderick-Westrope/teatime/server/internal/db"
)

func (app *application) handleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var creds entity.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			app.log.ErrorContext(ctx, "failed to unmarshal signup request body", slog.Any("error", err))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = app.repo.CreateUser(creds.Username, creds.Password)
		if err != nil {
			app.writeInternalServerError(ctx, w, "failed to create user", err)
			return
		}

		err = app.addNewSessionID(ctx, w, creds.Username)
		if err != nil {
			app.writeInternalServerError(ctx, w, "failed to set session ID", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("User created"))
		if err != nil {
			app.log.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		}
	}
}

func (app *application) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var creds entity.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			app.log.ErrorContext(ctx, "failed to unmarshal login request body", slog.Any("error", err))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		isAuthenticated, err := app.repo.AuthenticateUser(creds.Username, creds.Password)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				app.log.DebugContext(ctx, "user failed authentication", slog.Any("error", err))
				http.Error(w, "Failed authentication", http.StatusUnauthorized)
				return
			}
			app.writeInternalServerError(ctx, w, "failed to perform user authentication", err)
			return
		}

		if !isAuthenticated {
			app.log.DebugContext(ctx, "user failed authentication")
			http.Error(w, "Failed authentication", http.StatusUnauthorized)
			return
		}

		err = app.addNewSessionID(r.Context(), w, creds.Username)
		if err != nil {
			app.writeInternalServerError(ctx, w, "failed to set session ID", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Logged in"))
		if err != nil {
			app.log.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		}
	}
}

func (app *application) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		username, ok := r.Context().Value(ctxKeyUsername).(string)
		if !ok {
			app.writeInternalServerError(ctx, w, "failed to cast username to string", nil)
			return
		}

		err := app.repo.DeleteUserSessions(r.Context(), username)
		if err != nil {
			app.writeInternalServerError(ctx, w, "failed to delete user sessions", err)
			return
		}

		app.deleteSessionID(w)

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte("Logged out"))
		if err != nil {
			app.log.ErrorContext(ctx, "failed to write response", slog.Any("error", err))
		}
	}
}

func (app *application) handleWebSocket(workCtx context.Context, wg *sync.WaitGroup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		wg.Add(1)
		defer wg.Done()

		username, ok := ctx.Value(ctxKeyUsername).(string)
		if !ok {
			app.writeInternalServerError(ctx, w, "failed to get username from request context", nil)
			return
		}

		conn, err := app.hub.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		err = app.hub.NotifyConnection(username, true)
		if err != nil {
			app.writeInternalServerError(ctx, w, "failed to notify connection", err)
			return
		}
		app.hub.Add(conn, username)
		app.log.InfoContext(ctx, "new client connected", slog.String("username", username))

		defer func() {
			app.hub.Remove(username)
			err = app.hub.NotifyConnection(username, false)
			if err != nil {
				app.writeInternalServerError(ctx, w, "failed to notify disconnection", err)
				return
			}
			app.log.InfoContext(ctx, "client disconnected", slog.String("username", username))
		}()

		app.processWebSocketMessages(ctx, workCtx, conn, username)
	}
}

func (app *application) processWebSocketMessages(ctx, workCtx context.Context, conn *gws.Conn, username string) {
	msgCh := make(chan []byte)
	errCh := make(chan error)
	go func() {
		for {
			var msgData []byte
			_, msgData, err := conn.ReadMessage()
			if err != nil {
				errCh <- err
				return
			}
			msgCh <- msgData
		}
	}()

	for {
		select {
		case <-workCtx.Done():
			app.log.InfoContext(ctx, "closing connection", slog.String("username", username))
			err := app.hub.Close(username)
			if err != nil {
				app.log.ErrorContext(ctx, "failed to send close message",
					slog.String("username", username), slog.Any("error", err))
			}
			return

		case err := <-errCh:
			if websocket.IsNormalCloseError(err) {
				app.log.InfoContext(ctx, "received close message",
					slog.String("username", username), slog.String("value", err.Error()))
				return
			}
			app.log.ErrorContext(ctx, "error reading message",
				slog.String("username", username), slog.Any("error", err))
			return

		case msgData := <-msgCh:
			var msg websocket.Msg
			if err := json.Unmarshal(msgData, &msg); err != nil {
				app.log.ErrorContext(ctx, "error unmarshalling JSON",
					slog.String("username", username), slog.Any("error", err))
				return
			}

			switch payload := msg.Payload.(type) {
			case websocket.PayloadSendChatMessage:
				app.log.DebugContext(ctx, "sending message",
					slog.Any("recipients", payload.Recipients))
				err := app.hub.Send(msgData, payload.Recipients)
				if err != nil {
					app.log.ErrorContext(ctx, "failed to send message",
						slog.String("username", username), slog.Any("error", err))
				}

			default:
				app.log.ErrorContext(ctx, "message type has no handler",
					slog.String("username", username), slog.Any("msg", msg))
				return
			}
		}
	}
}
