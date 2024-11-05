package modals

import (
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
)

const (
	formKeyConversationName   = "conversation_name"
	formKeyParticipants       = "participants"
	formKeyNotifyParticipants = "notify_participants"
)

var _ tui.Modal = &CreateConversationModel{}

type CreateConversationModel struct {
	form                   *huh.Form
	hasAnnouncedCompletion bool
}

func NewCreateConversationModel() *CreateConversationModel {
	return &CreateConversationModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Key(formKeyConversationName).
					Title("Conversation Name").
					CharLimit(100).
					Validate(func(s string) error {
						if s == "" {
							return errors.New("empty conversation name not allowed")
						}
						return nil
					}),
				huh.NewText().Key(formKeyParticipants).
					Title("Participants").
					Description("Enter each participants username on a separate line.").
					Validate(func(s string) error {
						if s == "" {
							return errors.New("at least one participant is required")
						}
						return nil
					}),
				huh.NewConfirm().Key(formKeyNotifyParticipants).
					Title("Notify Participants?").
					Description("This will send a message on your behalf."),
			),
		),
	}
}

func (m *CreateConversationModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *CreateConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		switch m.hasAnnouncedCompletion {
		case false:
			cmds = append(cmds, m.announceCompletion())
		default:
			return m, nil
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *CreateConversationModel) announceCompletion() tea.Cmd {
	m.hasAnnouncedCompletion = true

	chatName := m.form.GetString(formKeyConversationName)
	participants := strings.Split(m.form.GetString(formKeyParticipants), "\n")
	notifyParticipants := m.form.GetBool(formKeyNotifyParticipants)

	return tea.Batch(
		tui.CreateConversationCmd(chatName, participants, notifyParticipants),
		tui.CloseModalCmd,
	)
}

func (m *CreateConversationModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		"Create a New Conversation:\n",
		m.form.View(),
	)
}

func (m *CreateConversationModel) SetSize(width, _ int) {
	m.form = m.form.WithWidth(width)
}
