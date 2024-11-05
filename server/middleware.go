package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
)

type ctxKey int

const (
	ctxKeyUsername ctxKey = iota
)

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
