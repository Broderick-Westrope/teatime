package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"

	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"

	"github.com/Broderick-Westrope/teatime/client/internal/db"
	"github.com/Broderick-Westrope/teatime/client/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
)

func main() {
	os.Exit(run())
}

func run() int {
	app, err := newApp()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create application: %v\n", err)
		return 1
	}
	defer app.cleanupFunc()

	err = app.runTui()
	if err != nil {
		app.log.Error("failed to run tui", slog.Any("error", err))
		return 1
	}
	return 0
}

type application struct {
	serverAddr string

	log         *slog.Logger
	debugID     string
	msgCh       chan tea.Msg
	cleanupFunc func()
}

func newApp() (*application, error) {
	app := &application{
		msgCh: make(chan tea.Msg),
	}

	err := app.loadEnvVars()
	if err != nil {
		return nil, err
	}

	var logFile *os.File
	if app.debugID != "" {
		logFile, err = createFilepath(fmt.Sprintf("logs/client_logs-%s.log", sanitizePathString(app.debugID)))
		if err != nil {
			return nil, err
		}
	}

	app.log = slog.New(slog.NewTextHandler(logFile, nil))
	app.cleanupFunc = func() {
		if logFile != nil {
			logFile.Close()
		}
	}
	return app, err
}

func (app *application) loadEnvVars() error {
	err := godotenv.Load(".client.env")
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .client.env file: %w", err)
	}
	app.debugID = os.Getenv("DEBUG")

	var found bool
	if app.serverAddr, found = os.LookupEnv("SERVER_ADDR"); !found {
		return fmt.Errorf("SERVER_ADDR env variable is required")
	}
	return nil
}

func (app *application) runTui() error {
	defer close(app.msgCh)

	var messagesDump *os.File
	var err error
	if app.debugID != "" {
		messagesDump, err = createFilepath(fmt.Sprintf("logs/client_messages-%s.log", app.debugID))
		if err != nil {
			return fmt.Errorf("failed to setup messages dump file: %w", err)
		}
		defer messagesDump.Close()
	}

	databaseFilePath, err := setupDatabaseFile()
	if err != nil {
		return fmt.Errorf("failed to setup database file: %w", err)
	}
	repo, err := db.NewRepository(fmt.Sprintf("file:%s", databaseFilePath))
	if err != nil {
		return fmt.Errorf("failed to setup database repository: %w", err)
	}

	m, err := starter.NewModel(app.msgCh, app.serverAddr, repo, messagesDump)
	if err != nil {
		return fmt.Errorf("failed to create starter model: %w", err)
	}
	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	go func() {
		for msg := range app.msgCh {
			p.Send(msg)
		}
	}()

	exitModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	typedExitModel, ok := exitModel.(*starter.Model)
	if !ok {
		return fmt.Errorf("failed to assert starter model type: %w", err)
	}

	if err = typedExitModel.ExitError; err != nil {
		if !websocket.IsNormalCloseError(err) {
			return fmt.Errorf("starter model exited with an error: %w", err)
		}
		app.log.Info("server disconnected gracefully", slog.Any("error", err))
	}
	return nil
}

func createFilepath(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}

		dir := filepath.Dir(path)
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			return nil, fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file '%s': %w", path, err)
	}
	return file, nil
}

func setupDatabaseFile() (string, error) {
	path, err := xdg.DataFile("TeaTime/client.db")
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

func sanitizePathString(s string) string {
	invalidChars := regexp.MustCompile(`[<>:"/\\|?* ]`)
	return invalidChars.ReplaceAllString(s, "-")
}
