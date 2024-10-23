package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	"github.com/Broderick-Westrope/teatime/internal/tui/views"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	messagesDumpFilepath = "logs/messages.log"
)

func main() {
	//setTestData()

	var messagesDump *os.File
	var err error
	if _, ok := os.LookupEnv("DEBUG"); ok {
		messagesDump, err = createFilepath(messagesDumpFilepath)
	}

	contacts := getTestData()

	m := starter.NewModel(
		views.NewAppModel(contacts, "Cordia_Tromp"),
		messagesDump,
	)

	opts := []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}

	exitModel, err := tea.NewProgram(m, opts...).Run()
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

func getTestData() []data.Contact {
	b, err := os.ReadFile("testdata.json")
	if err != nil {
		panic("failed to read testdata file: " + err.Error())
	}

	var contacts []data.Contact
	err = json.Unmarshal(b, &contacts)
	if err != nil {
		panic("failed to unmarshal testdata: " + err.Error())
	}

	return contacts
}

func setTestData() {
	time1, _ := time.Parse(time.RFC1123, "Sun, 12 Dec 2021 12:23:00 UTC")
	time2, _ := time.Parse(time.RFC1123, "Sun, 13 Dec 2021 12:23:00 UTC")
	contacts := []data.Contact{
		{
			Username: "Maynard.Adams",
			Conversation: []data.Message{
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
			Username: "Sherwood27",
			Conversation: []data.Message{
				{
					Content: "provident nesciunt sit",
				},
			},
		},
		{
			Username: "Elda48",
			Conversation: []data.Message{
				{
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed. Error quisquam rem ab.\nIpsum nam quasi exercitationem.\nMagni harum ipsum sit.\nA odit iusto provident.\nEaque eveniet tenetur porro tempora sint aut labore qui ea.",
				},
				{
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
