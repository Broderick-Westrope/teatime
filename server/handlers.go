package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Broderick-Westrope/teatime/internal/websocket"
)

func (app *application) handleWebSocket(ctx context.Context, wg *sync.WaitGroup) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()

		username := r.Header.Get("username")
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
