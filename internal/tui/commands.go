package tui

import tea "github.com/charmbracelet/bubbletea"

// FatalErrorCmd returns a command for creating a new FatalErrorMsg with the given error.
func FatalErrorCmd(err error) tea.Cmd {
	return func() tea.Msg {
		return FatalErrorMsg(err)
	}
}

// UpdateChatCmd is a command for creating a new UpdateChatMsg.
func UpdateChatCmd() tea.Msg {
	return UpdateChatMsg{}
}

func DebugLogCmd(value string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg(value)
	}
}
