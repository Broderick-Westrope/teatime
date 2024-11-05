package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const cookieNameSessionID = "session_id"

func (app *application) addNewSessionID(ctx context.Context, w http.ResponseWriter, username string) error {
	sessionID, err := app.repo.GetNewSessionID(ctx, username)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:    cookieNameSessionID,
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
	})
	return err
}

func (app *application) deleteSessionID(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieNameSessionID,
		Value:  "",
		MaxAge: -1,
	})
}

func (app *application) writeInternalServerError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	app.log.ErrorContext(ctx, msg, slog.Any("error", err))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
