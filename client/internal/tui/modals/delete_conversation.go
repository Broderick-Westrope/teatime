package modals

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
	"github.com/Broderick-Westrope/teatime/internal/entity"
)

const (
	formKeyConfirmDelete = "confirmDelete"
)

var _ tui.Modal = &DeleteConversationModel{}

type DeleteConversationModel struct {
	form                   *huh.Form
	hasAnnouncedCompletion bool
	conversationMD         entity.ConversationMetadata
}

func NewDeleteConversationModel(conversationMD entity.ConversationMetadata) *DeleteConversationModel {
	return &DeleteConversationModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Key(formKeyConfirmDelete),
			),
		),
		conversationMD: conversationMD,
	}
}

func (m *DeleteConversationModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *DeleteConversationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *DeleteConversationModel) announceCompletion() tea.Cmd {
	m.hasAnnouncedCompletion = true
	cmds := []tea.Cmd{tui.CloseModalCmd}

	if m.form.GetBool(formKeyConfirmDelete) {
		cmds = append(cmds, tui.DeleteConversationCmd(m.conversationMD))
	}

	return tea.Batch(cmds...)
}

func (m *DeleteConversationModel) View() string {
	return lipgloss.JoinVertical(lipgloss.Center,
		fmt.Sprintf("Are you sure you want to delete %q?\n", ansi.Truncate(m.conversationMD.Name, 35, "…")),
		m.form.View(),
	)
}

func (m *DeleteConversationModel) SetSize(width, _ int) {
	m.form = m.form.WithWidth(width)
}
