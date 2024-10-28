package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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
	ctx, cancelCtx := context.WithCancel(context.Background())
	app := newApp()
	server := &http.Server{Addr: serverAddress}
	wg := &sync.WaitGroup{}

	http.HandleFunc("/ws", app.handleWebSocket(ctx, wg))

	go app.startServer(server)
	app.handleShutdown(server, cancelCtx, wg)
}

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

func (app *application) startServer(server *http.Server) {
	app.log.Info("starting server", slog.String("addr", server.Addr))
	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		app.log.Error("HTTP server failed", slog.Any("error", err))
	}
	app.log.Info("stopped serving new connections", slog.Any("error", err))
}

func (app *application) handleShutdown(server *http.Server, cancelCtx context.CancelFunc, wg *sync.WaitGroup) {
	shutdownSigCh := make(chan os.Signal, 1)
	signal.Notify(shutdownSigCh, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSigCh

	app.log.Info("shutdown signal received")
	cancelCtx()
	wg.Wait()

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		app.log.Error("failed to shutdown HTTP server", slog.Any("error", err))
		os.Exit(1)
	}
	app.log.Info("graceful shutdown complete")
}
