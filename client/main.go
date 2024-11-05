package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Broderick-Westrope/teatime/client/internal/db"
	"github.com/Broderick-Westrope/teatime/client/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
)

type application struct {
	serverAddr string

	log     *slog.Logger
	isDebug bool
	msgCh   chan tea.Msg
}

func newApp(logWriter io.Writer) *application {
	_, isDebug := os.LookupEnv("DEBUG")

	return &application{
		serverAddr: "http://localhost:8080/",

		log:     slog.New(slog.NewTextHandler(logWriter, nil)),
		isDebug: isDebug,
		msgCh:   make(chan tea.Msg),
	}
}

func main() {
	os.Exit(run())
}

func run() int {
	logFile, err := createFilepath(fmt.Sprintf("logs/client_logs-%s.log", sanitizePathString(time.Now().String())))
	if err != nil {
		return 1
	}
	defer logFile.Close()

	app := newApp(logFile)

	err = app.runTui()
	if err != nil {
		app.log.Error("failed to run tui", slog.Any("error", err))
		return 1
	}
	return 0
}

func (app *application) runTui() error {
	defer close(app.msgCh)

	var messagesDump *os.File
	var err error
	if app.isDebug {
		messagesDump, err = createFilepath("logs/client-messages.log")
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
