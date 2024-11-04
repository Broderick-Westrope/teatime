package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/go-chi/chi/v5"
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
	r := chi.NewRouter()
	wg := &sync.WaitGroup{}

	r.Get("/ws", app.handleWebSocket(ctx, wg))

	server := &http.Server{
		Addr:    serverAddress,
		Handler: r,
	}
	go app.startServer(server)
	app.handleShutdown(server, cancelCtx, wg)
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
