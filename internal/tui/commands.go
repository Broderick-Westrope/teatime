package tui

import tea "github.com/charmbracelet/bubbletea"

func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}
