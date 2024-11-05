package views

import (
	"fmt"

	"github.com/Broderick-Westrope/teatime/client/internal/tui"
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
	errMessage             string

	width  int
	height int
}

func NewLockModel(errMessage string) *LockModel {
	return &LockModel{
		form: huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().Key(formKey_authMode).
					Affirmative("Signup").Negative("Login"),
				huh.NewInput().Key(formKey_username).
					Title("Username:").CharLimit(100).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("empty username not allowed")
						}
						return nil
					}),
				huh.NewInput().Key(formKey_password).
					Title("Password:").CharLimit(100).
					EchoMode(huh.EchoModePassword).
					Validate(func(s string) error {
						if s == "" {
							return fmt.Errorf("empty password not allowed")
						}
						return nil
					}),
			),
		),
		styles:     defaultLockStyles(),
		errMessage: errMessage,
	}
}

func (m *LockModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *LockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		formWidth := msg.Width - (lipgloss.Width(logoStr) + m.styles.Logo.GetHorizontalFrameSize() + m.styles.Form.GetHorizontalFrameSize())
		formWidth = min(formWidth, 50)
		m.form = m.form.WithWidth(formWidth)
		return m, nil
	}

	if m.form.State == huh.StateCompleted {
		switch m.hasAnnouncedCompletion {
		case false:
			return m, m.announceCompletion()
		default:
			return m, nil
		}
	}

	var cmds []tea.Cmd
	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *LockModel) announceCompletion() tea.Cmd {
	isSignup := m.form.GetBool(formKey_authMode)
	username := m.form.GetString(formKey_username)
	password := m.form.GetString(formKey_password)

	m.hasAnnouncedCompletion = true
	return tui.AuthenticateCmd(isSignup, username, password)
}

func (m *LockModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		m.styles.Title.Render(titleStr),
		creditStr,
		m.styles.ErrMessage.Render(m.errMessage),
		lipgloss.JoinHorizontal(lipgloss.Center,
			m.styles.Logo.Render(logoStr),
			m.styles.Form.Render(m.form.View()),
		),
	)
	output = m.styles.View.Render(output)
	return lipgloss.Place(m.width, (m.height/10)*9, lipgloss.Center, lipgloss.Center, output)
}
