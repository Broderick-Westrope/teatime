package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/Broderick-Westrope/teatime/server/internal/db"
)

func main() {
	os.Exit(run())
}

func run() int {
	app, err := newApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create application: %s", err)
		return 1
	}

	ctx, cancelCtx := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	server := &http.Server{
		Addr:              app.serverAddr,
		Handler:           app.setupRouter(ctx, wg),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1 MB
	}

	go app.startServer(server)
	app.handleShutdown(server, cancelCtx, wg)
	return 0
}

type application struct {
	hub  *websocket.Hub
	log  *slog.Logger
	repo *db.Repository

	serverAddr string
	redisAddr  string
	dbConn     string
	logLevel   slog.Level
}

func newApp() (*application, error) {
	app := &application{
		hub: websocket.NewHub(),
	}
	err := app.loadEnvVars()
	if err != nil {
		return nil, err
	}

	app.repo, err = db.NewRepository(app.dbConn, app.redisAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database repository: %w", err)
	}
	app.log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: app.logLevel}))

	return app, nil
}

func (app *application) loadEnvVars() error {
	err := godotenv.Load(".server.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .server.env file: %w", err)
	}

	var found bool
	if app.serverAddr, found = os.LookupEnv("SERVER_ADDR"); !found {
		return fmt.Errorf("SERVER_ADDR env variable is required")
	}
	if app.redisAddr, found = os.LookupEnv("REDIS_ADDR"); !found {
		return fmt.Errorf("REDIS_ADDR env variable is required")
	}
	if app.dbConn, found = os.LookupEnv("DB_CONN"); !found {
		return fmt.Errorf("DB_CONN env variable is required")
	}
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "":
		app.logLevel = slog.LevelInfo
	default:
		var intLevel int
		intLevel, err = strconv.Atoi(logLevel)
		if err != nil {
			return fmt.Errorf("LOG_LEVEL env variable must be an integer")
		}
		app.logLevel = slog.Level(intLevel)
	}
	return nil
}

func (app *application) setupRouter(ctx context.Context, wg *sync.WaitGroup) chi.Router {
	r := chi.NewRouter().With(app.loggerMiddleware(), middleware.Recoverer)

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
		return
	}
	app.log.Info("graceful shutdown complete")
}
