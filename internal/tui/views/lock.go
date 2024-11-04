package views

import (
	"github.com/Broderick-Westrope/teatime/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

const (
	titleStr = `
 _______      _______
|__   __|    |__   __|
   | | ___  __ _| (_)_ __ ___   ___
   | |/ _ \/ _' | | | '_ ' _ \ / _ \
   | |  __/ (_| | | | | | | | |  __/
   |_|\___|\__,_|_|_|_| |_| |_|\___|`
	creditStr = "Created by Broderick Westrope\nSource code at github.com/Broderick-Westrope/teatime"
	logoStr   = `
      {
   {   }
    }_{ __{
 .-{   }   }-.
(   }     {   )
|'-.._____..-'|
|             ;--.
|            (__  \
|             | )  )
|             |/  /
|             /  /
|            (  /
\             y'
 '-.._____..-'
  -Felix Lee-`
	formKey_authMode = "authMode"
	formKey_username = "username"
	formKey_password = "password"
)

var _ tea.Model = &LockModel{}

type LockModel struct {
	form                   *huh.Form
	hasAnnouncedCompletion bool
	styles                 *lockStyles
}

func NewLockModel() *LockModel {
	return &LockModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Key(formKey_authMode).
					Affirmative("Signup").Negative("Login"),
				huh.NewInput().Key(formKey_username).
					Title("Username:").CharLimit(100),
				huh.NewInput().Key(formKey_password).
					Title("Password:").CharLimit(100).
					EchoMode(huh.EchoModePassword),
			),
		),
		styles: defaultLockStyles(),
	}
}

func (m *LockModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *LockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		formWidth := msg.Width - (lipgloss.Width(logoStr) + m.styles.GetFormRelativeHorizontalFrameSize())
		formWidth = min(formWidth, 50)
		m.form = m.form.WithWidth(formWidth) //.WithHeight(formHeight)
		//if lipgloss.Height(logoStr) > lipgloss.Height(m.form.View()) {
		//	m.form = m.form.WithHeight(lipgloss.Height(logoStr))
		//}
		m.styles.View = m.styles.View.Width(msg.Width).Height(msg.Height)
		return m, nil
	}

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

func (m *LockModel) announceCompletion() tea.Cmd {
	m.hasAnnouncedCompletion = true

	return tui.AuthenticateCmd(
		m.form.GetBool(formKey_authMode),
		m.form.GetString(formKey_username),
		m.form.GetString(formKey_password),
	)
}

func (m *LockModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		m.styles.Title.Render(titleStr),
		creditStr,
		lipgloss.JoinHorizontal(lipgloss.Center,
			m.styles.Logo.Render(logoStr),
			m.styles.Form.Render(m.form.View()),
		),
	)
	return m.styles.View.Render(output)
}
