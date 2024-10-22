package views

import "github.com/charmbracelet/lipgloss"

type AppStyles struct {
	View     lipgloss.Style
	Contacts lipgloss.Style
	Chat     lipgloss.Style
}

func DefaultAppStyles() *AppStyles {
	return &AppStyles{
		View:     lipgloss.NewStyle().Margin(1).Align(lipgloss.Center, lipgloss.Center),
		Contacts: lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(lipgloss.Color("237")).PaddingRight(2),
		Chat:     lipgloss.NewStyle(),
	}
}
