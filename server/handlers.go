package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/Broderick-Westrope/teatime/server/internal/db"
)

func (app *application) handleSignup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds entity.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			app.log.Error("failed to unmarshal signup request body", slog.Any("error", err))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err = app.repo.CreateUser(creds.Username, creds.Password)
		if err != nil {
			app.writeInternalServerError(w, "failed to create user", err)
			return
		}

		err = app.addNewSessionID(r.Context(), w, creds.Username)
		if err != nil {
			app.writeInternalServerError(w, "failed to set session ID", err)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created"))
	}
}

func (app *application) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds entity.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			app.log.Error("failed to unmarshal login request body", slog.Any("error", err))
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		isAuthenticated, err := app.repo.AuthenticateUser(creds.Username, creds.Password)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				app.log.Debug("user failed authentication", slog.Any("error", err))
				http.Error(w, "Failed authentication", http.StatusUnauthorized)
				return
			}
			app.writeInternalServerError(w, "failed to perform user authentication", err)
			return
		}

		if !isAuthenticated {
			app.log.Debug("user failed authentication")
			http.Error(w, "Failed authentication", http.StatusUnauthorized)
			return
		}

		err = app.addNewSessionID(r.Context(), w, creds.Username)
		if err != nil {
			app.writeInternalServerError(w, "failed to set session ID", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged in"))
	}
}

func (app *application) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, ok := r.Context().Value("username").(string)
		if !ok {
			app.writeInternalServerError(w, "failed to cast username to string", nil)
			return
		}

		err := app.repo.DeleteUserSessions(r.Context(), username)
		if err != nil {
			app.writeInternalServerError(w, "failed to delete user sessions", err)
			return
		}

		app.deleteSessionID(w)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged in"))
	}
}

func (app *application) handleWebSocket(ctx context.Context, wg *sync.WaitGroup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()

		username, ok := r.Context().Value("username").(string)
		if !ok {
			app.writeInternalServerError(w, "failed to get username from request context", nil)
			return
		}
		usernameAttr := slog.String("username", username)

		conn, err := app.hub.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		app.hub.Add(conn, username)
		app.log.Info("new client connected", usernameAttr)
		defer func() {
			app.hub.Remove(username)
			app.log.Info("client disconnected", usernameAttr)
		}()

		msgCh := make(chan []byte)
		errCh := make(chan error)
		go func() {
			for {
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
			case <-ctx.Done():
				app.log.Info("closing connection", usernameAttr)
				err = app.hub.Close(username)
				if err != nil {
					app.log.Error("failed to send close message", usernameAttr, slog.Any("error", err))
				}
				return

			case err = <-errCh:
				if websocket.IsNormalCloseError(err) {
					app.log.Info("received close message", usernameAttr, slog.String("value", err.Error()))
					return
				}
				app.log.Error("error reading message", usernameAttr, slog.Any("error", err))
				return

			case msgData := <-msgCh:
				var msg websocket.Msg
				if err = json.Unmarshal(msgData, &msg); err != nil {
					app.log.Error("error unmarshalling JSON", usernameAttr, slog.Any("error", err))
					return
				}

				switch payload := msg.Payload.(type) {
				case websocket.PayloadSendChatMessage:
					app.log.Debug("sending message", slog.Any("recipients", payload.Recipients))
					err = app.hub.Send(msgData, payload.Recipients)
					if err != nil {
						app.log.Error("failed to send message", usernameAttr, slog.Any("error", err))
					}

				default:
					app.log.Error("message type has no handler", usernameAttr, slog.Any("msg", msg))
					return
				}
			}
		}
	}
}
