package views

import (
	"errors"
	"fmt"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/components"
	"github.com/Broderick-Westrope/teatime/internal/tui/components/contacts"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	ErrInvalidTypeAssertion = errors.New("invalid type assertion")
)

type AppModel struct {
	contacts      *contacts.Model
	chat          *components.ChatModel
	chatIsFocused bool
	styles        *AppStyles
}

func NewAppModel() *AppModel {
	contactItems := []contacts.Contact{
		{
			Username: "Maynard.Adams",
			Conversation: []data.Message{
				{
					Content: "Doloribus eligendi at velit qui.",
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
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed.",
				},
			},
		},
	}

	return &AppModel{
		contacts:      contacts.NewModel(contactItems),
		chat:          components.NewChatModel(contactItems[0].Conversation, contactItems[0].Username),
		chatIsFocused: false,
		styles:        DefaultAppStyles(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	return nil
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		frameX, frameY := m.styles.View.GetFrameSize()
		cmd, err := m.updateComponentSizes(msg.Width-frameX, msg.Height-frameY)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd
	}

	switch m.chatIsFocused {
	case true:
		cmd, err := m.updateChatModel(msg)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd

	case false:
		cmd, err := m.updateContactsModel(msg)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd
	}

	return m, nil
}

func (m *AppModel) View() string {
	var output string
	output = lipgloss.JoinHorizontal(lipgloss.Top,
		m.styles.Contacts.Render(m.contacts.View()),
		m.styles.Chat.Render(m.chat.View()),
	)
	return m.styles.View.Render(output)
}

func (m *AppModel) updateComponentSizes(width, height int) (tea.Cmd, error) {
	var cmds []tea.Cmd

	cmd, err := m.updateContactsModel(tui.ComponentSizeMsg{
		Width:  width / 3,
		Height: height,
	})
	if err != nil {
		return nil, err
	}
	cmds = append(cmds, cmd)

	cmd, err = m.updateChatModel(tui.ComponentSizeMsg{
		Width:  (width / 3) * 2,
		Height: height,
	})
	if err != nil {
		return nil, err
	}
	cmds = append(cmds, cmd)

	return tea.Batch(cmd), nil
}

func (m *AppModel) updateChatModel(msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := m.chat.Update(msg)
	m.chat, ok = newModel.(*components.ChatModel)
	if !ok {
		return nil, fmt.Errorf("failed to update chat model: %w", ErrInvalidTypeAssertion)
	}
	return cmd, nil
}

func (m *AppModel) updateContactsModel(msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := m.contacts.Update(msg)
	m.contacts, ok = newModel.(*contacts.Model)
	if !ok {
		return nil, fmt.Errorf("failed to update contacts model: %w", ErrInvalidTypeAssertion)
	}
	return cmd, nil
}
