package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"

	"github.com/Broderick-Westrope/teatime/internal/websocket"
)

const serverAddress = ":8080"

type application struct {
	hub *websocket.Hub
	log *slog.Logger
}

func newApp() *application {
	return &application{
		hub: websocket.NewHub(),
		log: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func main() {
	ctx := context.Background()
	app := newApp()

	http.HandleFunc("/ws", app.handleWebSocket(ctx))

	app.log.Info("starting server", slog.String("addr", serverAddress))
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		app.log.Error("failed to listen", slog.Any("error", err))
	}
}

func (app *application) handleWebSocket(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		for {
			select {
			case <-ctx.Done():
				return

			default:
				_, msgData, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsNormalCloseError(err) {
						app.log.Info("received close message", usernameAttr, slog.String("value", err.Error()))
						return
					}

					app.log.Error("error reading message", usernameAttr, slog.Any("error", err))
					return
				}

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
