package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

const serverAddress = ":8080"

type application struct {
	clientManager ClientManager
	log           *slog.Logger
}

func main() {
	ctx := context.Background()

	app := application{
		clientManager: NewClientManager(),
		log:           slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	http.HandleFunc("/ws", app.handleWebSocket(ctx))

	go app.broadcastTime(ctx)

	app.log.InfoContext(ctx, "starting server", slog.String("addr", serverAddress))
	err := http.ListenAndServe(serverAddress, nil)
	if err != nil {
		app.log.Error("failed to listen", slog.Any("error", err))
	}
}

func (app *application) handleWebSocket(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := app.clientManager.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		app.clientManager.Add(conn)
		app.log.InfoContext(ctx, "new client connected")
		defer func() {
			app.clientManager.Remove(conn)
			app.log.InfoContext(ctx, "client disconnected")
		}()

		for {
			select {
			case <-ctx.Done():
				return

			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
						app.log.InfoContext(ctx, "received close message", slog.String("value", err.Error()))
						return
					}

					app.log.ErrorContext(ctx, "error reading message", slog.Any("error", err))
					return
				}

				app.log.InfoContext(ctx, "received message", slog.String("string", string(message)))
			}
		}
	}
}

func (app *application) broadcastTime(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 4)

	for range ticker.C {
		msg := "the time is " + time.Now().String()

		err := app.clientManager.Broadcast([]byte(msg))
		if err != nil {
			slog.ErrorContext(ctx, "error broadcasting message", slog.Any("error", err))
		}
	}
}
