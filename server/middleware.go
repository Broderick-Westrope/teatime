package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type ctxKey int

const (
	ctxKeyUsername ctxKey = iota
)

// Auth ------------------------------

func (app *application) authMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					app.log.Debug("session ID cookie not found", slog.Any("error", err))
					http.Error(w, "Session ID is required", http.StatusBadRequest)
					return
				}

				app.log.Debug("failed to get session ID cookie", slog.Any("error", err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			username, err := app.repo.GetUsernameWithSessionID(r.Context(), cookie.Value)
			if err != nil {
				app.log.Debug("failed to get username with session ID", slog.Any("error", err))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUsername, username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Logging ------------------------------

type slogFormatter struct {
	Logger *slog.Logger
}

func (f *slogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &slogEntry{
		Logger:  f.Logger,
		Request: r,
	}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		entry.Logger = entry.Logger.With("request_id", reqID)
	}
	return entry
}

type slogEntry struct {
	Logger  *slog.Logger
	Request *http.Request
}

func (e *slogEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.Logger.With(
		"status", status,
		"bytes", bytes,
		"elapsed", elapsed,
		"method", e.Request.Method,
		"path", e.Request.URL.Path,
		"remote_addr", e.Request.RemoteAddr,
		"user_agent", e.Request.UserAgent(),
	).Info("handled request")
}

func (e *slogEntry) Panic(v interface{}, stack []byte) {
	e.Logger.With(
		"panic", v,
		"stack", string(stack),
	).Error("panic recovered")
}

func (app *application) loggerMiddleware() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&slogFormatter{Logger: app.log})
}
