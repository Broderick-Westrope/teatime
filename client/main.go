package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/tui/views"
	"github.com/Broderick-Westrope/teatime/internal/websocket"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//setTestData()

	username := os.Args[1]

	var messagesDump *os.File
	var err error
	if _, ok := os.LookupEnv("DEBUG"); ok {
		messagesDump, err = createFilepath(fmt.Sprintf("logs/messages_%s.log", sanitizePathString(username)))
	}

	contacts := getTestData()

	wsClient, err := websocket.NewClient("ws://localhost:8080/ws", username)
	if err != nil {
		log.Fatalf("failed to create WebSocket client: %v\n", err)
	}
	defer wsClient.Close()

	msgCh := make(chan tea.Msg)
	go readFromWebSocket(wsClient, msgCh)

	m := starter.NewModel(
		views.NewAppModel(contacts, username),
		wsClient,
		messagesDump,
	)

	opts := []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}

	p := tea.NewProgram(m, opts...)

	go func() {
		for msg := range msgCh {
			p.Send(msg)
		}
	}()

	exitModel, err := p.Run()
	if err != nil {
		log.Fatalf("alas, there's been an error: %v\n", err)
	}

	typedExitModel, ok := exitModel.(*starter.Model)
	if !ok {
		log.Fatalln("failed to assert starter model type")
	}

	if typedExitModel.ExitError != nil {
		log.Fatalf("starter model exited with an error: %v\n", typedExitModel.ExitError)
	}
}

func readFromWebSocket(client *websocket.Client, msgCh chan tea.Msg) {
	defer close(msgCh)

	for {
		msg, err := client.ReadMessage()
		if err != nil {
			panic(err)
			if websocket.IsNormalCloseError(err) {
				//app.log.InfoContext(ctx, "received close message", slog.String("value", err.Error()))
				return
			}
			//app.log.ErrorContext(ctx, "failed to read message", slog.Any("error", err))
			return
		}

		switch payload := msg.Payload.(type) {
		case websocket.PayloadSendChatMessage:
			msgCh <- tui.ReceiveMessageMsg{
				ConversationName: payload.ChatName,
				Message:          payload.Message,
			}
		default:
			panic("unknown payload")
		}

		//app.log.InfoContext(ctx, "received message", slog.String("value", string(msg)))
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

func getTestData() []data.Conversation {
	b, err := os.ReadFile("testdata.json")
	if err != nil {
		panic("failed to read testdata file: " + err.Error())
	}

	var conversations []data.Conversation
	err = json.Unmarshal(b, &conversations)
	if err != nil {
		panic("failed to unmarshal testdata: " + err.Error())
	}

	return conversations
}

func setTestData() {
	time1, _ := time.Parse(time.RFC1123, "Sun, 12 Dec 2021 12:23:00 UTC")
	time2, _ := time.Parse(time.RFC1123, "Sun, 13 Dec 2021 12:23:00 UTC")
	contacts := []data.Conversation{
		{
			Name: "Maynard.Adams",
			Participants: []string{
				"Maynard.Adams",
				"Cordia_Tromp",
			},
			Messages: []data.Message{
				{
					Author:  "Maynard.Adams",
					Content: "Doloribus eligendi at velit qui.",
					SentAt:  time1,
				},
				{
					Author:  "Cordia_Tromp",
					Content: "Earum similique tempore. Ullam animi hic repudiandae. Amet id voluptas id error veritatis tenetur incidunt quidem nihil. Eius facere nostrum expedita eum.\nDucimus in temporibus non. Voluptatum enim odio cupiditate error est aspernatur eligendi. Ea iure tenetur nam. Nemo quo veritatis iusto maiores illum modi necessitatibus. Sunt minus ab.\nOfficia deserunt omnis velit aliquid facere sit. Vel rem atque. Veniam dolores corporis quasi sit deserunt minus molestias sunt.",
					SentAt:  time2,
				},
			},
		},
		{
			Name: "Sherwood27",
			Participants: []string{
				"Sherwood27",
			},
			Messages: []data.Message{
				{
					Author:  "Sherwood27",
					Content: "provident nesciunt sit",
				},
			},
		},
		{
			Name: "Elda48",
			Participants: []string{
				"Elda48",
				"Jay Bernhard",
			},
			Messages: []data.Message{
				{
					Author:  "Elda48",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
					Author:  "Elda48",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
					Author:  "Jay Bernhard",
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
			},
		},
	}

	b, err := json.Marshal(contacts)
	if err != nil {
		panic("failed to marshal testdata: " + err.Error())
	}

	err = os.WriteFile("testdata.json", b, 0700)
	if err != nil {
		panic("failed to write testdata file: " + err.Error())
	}
}

func sanitizePathString(username string) string {
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	return invalidChars.ReplaceAllString(username, "_")
}
