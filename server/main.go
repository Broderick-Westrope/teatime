package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/Broderick-Westrope/teatime/server/internal/db"
	"github.com/adrg/xdg"
	"github.com/go-chi/chi/v5"
)

const (
	serverAddress = ":8080"
	redisAddress  = "redis:6379"
)

type application struct {
	hub  *websocket.Hub
	log  *slog.Logger
	repo *db.Repository
}

func newApp() (*application, error) {
	databaseFilePath, err := setupDatabaseFile()
	if err != nil {
		return nil, fmt.Errorf("failed to setup database file: %w", err)
	}
	repo, err := db.NewRepository(fmt.Sprintf("file:%s", databaseFilePath), redisAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database repository: %w", err)
	}

	return &application{
		hub:  websocket.NewHub(),
		log:  slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
		repo: repo,
	}, nil
}

func main() {
	ctx, cancelCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	r := chi.NewRouter()

	app, err := newApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create application: %s", err)
		os.Exit(1)
	}

	r.Get("/ws", app.handleWebSocket(ctx, wg))

	server := &http.Server{
		Addr:    serverAddress,
		Handler: app.setupRouter(ctx, wg),
	}
	go app.startServer(server)
	app.handleShutdown(server, cancelCtx, wg)
}

func (app *application) setupRouter(ctx context.Context, wg *sync.WaitGroup) chi.Router {
	r := chi.NewRouter()

	r.Route("/auth", func(r chi.Router) {
		r.Get("/signup", app.handleSignup())
		r.Get("/login", app.handleLogin())
		r.With(app.authMiddleware()).Get("/logout", app.handleLogout())
	})

	r.With(app.authMiddleware()).Get("/ws", app.handleWebSocket(ctx, wg))

	return r
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

func setupDatabaseFile() (string, error) {
	path, err := xdg.DataFile("TeaTime/server.db")
	if err != nil {
		return "", err
	}

	_, err = os.Stat(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}

		_, err = os.Create(path)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}
