package main

import (
	"log"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui/components/contacts"
	"github.com/Broderick-Westrope/teatime/internal/tui/starter"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//msgs := []data.Message{
	//	{
	//		Content: "some",
	//		Author:  "me",
	//	},
	//	{
	//		Content: "thing",
	//		Author:  "other",
	//	},
	//}

	items := []contacts.Contact{
		{
			Name: "Maynard.Adams",
			Conversation: []data.Message{
				{
					Content: "Doloribus eligendi at velit qui.",
				},
			},
		},
		{
			Name: "Sherwood27",
			Conversation: []data.Message{
				{
					Content: "provident nesciunt sit",
				},
			},
		},
		{
			Name: "Elda48",
			Conversation: []data.Message{
				{
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed.",
				},
			},
		},
	}

	m := starter.NewModel(
		//components.NewChatModel(msgs, "me", 50),
		contacts.NewModel(items),
	)

	_, err := tea.NewProgram(m).Run()
	if err != nil {
		log.Fatal("alas, there's been an error")
	}
}
