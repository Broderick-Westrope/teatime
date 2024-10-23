package views

import (
	"fmt"
	"time"

	"github.com/Broderick-Westrope/teatime/internal/data"
	"github.com/Broderick-Westrope/teatime/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/tui/components"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var _ tea.Model = &AppModel{}

type AppModel struct {
	contacts *components.ContactsModel
	chat     *components.ChatModel
	focus    FocusRegion
	styles   *AppStyles
}

type FocusRegion int

const (
	FocusRegionContacts FocusRegion = iota
	FocusRegionChat
)

func NewAppModel() *AppModel {
	time1, _ := time.Parse(time.RFC1123, "Sun, 12 Dec 2021 12:23:00 UTC")
	time2, _ := time.Parse(time.RFC1123, "Sun, 13 Dec 2021 12:23:00 UTC")

	contactItems := []components.Contact{
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
					Content: "Nulla eaque molestias molestiae porro iusto. Laboriosam sequi laborum autem harum iste ex. Autem minus pariatur soluta voluptatum. Quis dolores cumque atque quisquam unde. Aliquid officia veritatis nihil voluptate dolorum. Delectus recusandae natus ratione animi.\nQuasi unde dolor modi est libero quo quam iste eum. Itaque facere dolore dignissimos placeat. Cumque magni quia reprehenderit voluptas sequi voluptatum reprehenderit.\nAsperiores dolorum eum animi tempora laudantium autem. Omnis quidem atque laboriosam maiores laudantium. Fuga possimus mollitia amet adipisci rerum. Excepturi blanditiis libero modi harum sed.",
				},
			},
		},
	}

	focus := FocusRegionContacts

	return &AppModel{
		contacts: components.NewContactsModel(contactItems, focus == FocusRegionContacts),
		chat:     components.NewChatModel(contactItems[0].Conversation, "Cordia_Tromp", contactItems[0].Username, focus == FocusRegionChat),
		focus:    focus,
		styles:   DefaultAppStyles(),
	}
}

func (m *AppModel) Init() tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	cmd = m.contacts.Init()
	cmds = append(cmds, cmd)

	cmd = m.chat.Init()
	cmds = append(cmds, cmd)

	return tea.Batch(cmds...)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		frameWidth, frameHeight := m.styles.TotalFrameSize()
		cmd, err := m.updateComponentSizes(msg.Width-frameWidth, msg.Height-frameHeight)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd

	case tui.SetConversationMsg:
		contact, err := m.contacts.GetSelectedContact()
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		m.chat.SetConversation(contact.Conversation, contact.Username)
		err = m.setFocus(FocusRegionChat)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, nil

	case tui.SendMessageMsg:
		cmd, err := m.contacts.AddNewMessage(msg)
		if err != nil {
			return m, tui.FatalErrorCmd(err)
		}
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// move from chat to contacts
			if m.focus != FocusRegionChat {
				break
			}
			err := m.setFocus(FocusRegionContacts)
			if err != nil {
				return m, tui.FatalErrorCmd(err)
			}
			m.chat.ResetInput()
			return m, nil

		case "q":
			// quit unless typing "q" in the chat
			if m.focus == FocusRegionChat {
				break
			}
			return m, tea.Quit
		}
	}

	cmd, err := m.updateFocussedChild(msg)
	if err != nil {
		return m, tui.FatalErrorCmd(err)
	}
	return m, cmd
}

func (m *AppModel) View() string {
	var output string
	output = lipgloss.JoinHorizontal(lipgloss.Center,
		m.styles.Contacts.Render(m.contacts.View()),
		m.styles.Chat.Render(m.chat.View()),
	)
	return m.styles.View.Render(output)
}

func (m *AppModel) updateFocussedChild(msg tea.Msg) (tea.Cmd, error) {
	switch m.focus {
	case FocusRegionContacts:
		return m.updateContactsModel(msg)

	case FocusRegionChat:
		return m.updateChatModel(msg)

	default:
		return nil, fmt.Errorf("unknown FocusRegion %d", m.focus)
	}
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
		return nil, fmt.Errorf("failed to update chat model: %w", tui.ErrInvalidTypeAssertion)
	}
	return cmd, nil
}

func (m *AppModel) updateContactsModel(msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := m.contacts.Update(msg)
	m.contacts, ok = newModel.(*components.ContactsModel)
	if !ok {
		return nil, fmt.Errorf("failed to update contacts model: %w", tui.ErrInvalidTypeAssertion)
	}
	return cmd, nil
}

func (m *AppModel) setFocus(focus FocusRegion) error {
	switch focus {
	case FocusRegionContacts:
		m.chat.SwitchStyleFunc(components.DisabledChatStyleFunc)
		m.contacts.SwitchStyles(components.EnabledContactsStyles())

	case FocusRegionChat:
		m.chat.SwitchStyleFunc(components.EnabledChatStyleFunc)
		m.contacts.SwitchStyles(components.DisabledContactsStyles())

	default:
		return fmt.Errorf("unknown FocusRegion %d", focus)
	}
	m.focus = focus
	return nil
}
