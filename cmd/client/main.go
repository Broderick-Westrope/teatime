package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/db"
	"github.com/Broderick-Westrope/teatime/internal/entity"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	"github.com/adrg/xdg"
	tea "github.com/charmbracelet/bubbletea"
)

type application struct {
	username string
	password string

	log      *slog.Logger
	wsClient *websocket.Client
	isDebug  bool
	msgCh    chan tea.Msg
}

func newApp(username, password string, logWriter io.Writer) *application {
	wsClient, err := websocket.NewClient("ws://localhost:8080/ws", username)
	if err != nil {
		log.Fatalf("failed to create WebSocket client: %v\n", err)
	}

	_, isDebug := os.LookupEnv("DEBUG")

	return &application{
		username: username,
		password: password,

		log:      slog.New(slog.NewTextHandler(logWriter, nil)),
		wsClient: wsClient,
		isDebug:  isDebug,
		msgCh:    make(chan tea.Msg),
	}
}

func main() {
	username := os.Args[1]
	password := os.Args[2]

	logFile, err := createFilepath(fmt.Sprintf("logs/client_%s.log", sanitizePathString(username)))
	if err != nil {
		panic(fmt.Sprintf("failed to create log file: %v", err))
	}
	defer logFile.Close()

	app := newApp(username, password, logFile)
	defer app.wsClient.Close()

	go app.readFromWebSocket()

	err = app.runTui()
	if err != nil {
		app.log.Error("failed to run tui", slog.Any("error", err))
		os.Exit(1)
	}
}

func (app *application) runTui() error {
	var messagesDump *os.File
	var err error
	if app.isDebug {
		messagesDump, err = createFilepath(fmt.Sprintf("logs/client-messages_%s.log", sanitizePathString(app.username)))
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

	m, err := starter.NewModel(app.username, app.password, app.wsClient, repo, messagesDump)
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

func (app *application) readFromWebSocket() {
	defer close(app.msgCh)

	for {
		msg, err := app.wsClient.ReadMessage()
		if err != nil {
			app.log.Error("failed to read message", slog.Any("error", err))
			app.msgCh <- tui.FatalErrorMsg(err)
			return
		}

		switch payload := msg.Payload.(type) {
		case websocket.PayloadSendChatMessage:
			app.msgCh <- tui.ReceiveMessageMsg{
				ConversationName: payload.ConversationName,
				Message:          payload.Message,
			}
		default:
			app.log.Error("unknown WebSocket message payload", slog.Int("msg_type", int(msg.Type)))
			app.msgCh <- tui.FatalErrorMsg(fmt.Errorf("unknown WebSocket message payload"))
			return
		}

		// TODO: This can be removed after debugging
		app.log.Info("received message", slog.Any("value", msg))
	}
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
			return nil, fmt.Errorf("failed to create directory '%s': %w\n", dir, err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file '%s': %w\n", path, err)
	}
	return file, nil
}

func setupDatabaseFile() (string, error) {
	path, err := xdg.DataFile("TeaTime/teatime.db")
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

func getTestData() []entity.Conversation {
	time1, _ := time.Parse(time.RFC1123, "Sun, 12 Dec 2021 12:23:00 UTC")
	time2, _ := time.Parse(time.RFC1123, "Sun, 13 Dec 2021 12:23:00 UTC")
	conversations := []entity.Conversation{
		{
			Name: "TEST CHAT",
			Participants: []string{
				"Sally.Sender",
				"Robby.Receiver",
				// With only two people the conversation name is the other persons name.
				// This is not possible with hardcoded test data. The solution is to add a third,
				// unused person so that the manually set conversation name will be used.
				"A third person",
			},
			Messages: []entity.Message{
				{
					Author:  "Sally.Sender",
					Content: "Doloribus eligendi at velit qui.",
					SentAt:  time1,
				},
				{
					Author:  "Robby.Receiver",
					Content: "Earum similique tempore. Ullam animi hic repudiandae. Amet id voluptas id error veritatis tenetur incidunt quidem nihil. Eius facere nostrum expedita eum.\nDucimus in temporibus non. Voluptatum enim odio cupiditate error est aspernatur eligendi. Ea iure tenetur nam. Nemo quo veritatis iusto maiores illum modi necessitatibus. Sunt minus ab.\nOfficia deserunt omnis velit aliquid facere sit. Vel rem atque. Veniam dolores corporis quasi sit deserunt minus molestias sunt.",
					SentAt:  time2,
				},
			},
		},
		{
			Name: "Sherwood27",
			Participants: []string{
				"Sherwood27",
				"Sally.Sender",
			},
			Messages: []entity.Message{
				{
					Author:  "Sherwood27",
					Content: "provident nesciunt sit",
				},
			},
		},
		{
			Name: "Rick48",
			Participants: []string{
				"Rick48",
				"Robby.Receiver",
			},
			Messages: []entity.Message{
				{
					Author:  "Rick48",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
					Author:  "Rick48",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
					Author:  "Robby.Receiver",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
			},
		},
	}

	return conversations
}

func sanitizePathString(username string) string {
	invalidChars := regexp.MustCompile(`[<>:"/\\|?* ]`)
	return invalidChars.ReplaceAllString(username, "-")
}
