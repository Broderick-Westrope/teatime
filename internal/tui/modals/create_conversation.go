package modals

import (
	"strings"

	"github.com/Broderick-Westrope/teatime/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	formKey_conversationName   = "conversation_name"
	formKey_participants       = "participants"
	formKey_notifyParticipants = "notify_participants"
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
				huh.NewInput().Key(formKey_conversationName).
					Title("Conversation Name").
					CharLimit(100),
				huh.NewText().Key(formKey_participants).
					Title("Participants").
					Description("Enter each participants username on a separate line."),
				huh.NewConfirm().Key(formKey_notifyParticipants).
					Title("Notify Participants?"),
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

	chatName := m.form.GetString(formKey_conversationName)
	participants := strings.Split(m.form.GetString(formKey_participants), "\n")
	notifyParticipants := m.form.GetBool(formKey_notifyParticipants)

	switch {
	case len(participants) <= 0:
		return tui.CloseModalCmd

	case chatName == "":
		switch {
		case len(participants) == 1:
			chatName = participants[0]
		default:
			return tui.CloseModalCmd
		}
	}

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

func (m *CreateConversationModel) SetSize(width, height int) {
	m.form = m.form.WithWidth(width) //.WithHeight(height)
}
