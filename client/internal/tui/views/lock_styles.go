package views

import "github.com/charmbracelet/lipgloss"

type lockStyles struct {
	Title      lipgloss.Style
	ErrMessage lipgloss.Style
	Logo       lipgloss.Style
	Form       lipgloss.Style
	View       lipgloss.Style
}

func defaultLockStyles() *lockStyles {
	return &lockStyles{
		Title:      lipgloss.NewStyle().Margin(1, 2).AlignHorizontal(lipgloss.Left),
		ErrMessage: lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		Logo:       lipgloss.NewStyle().Margin(1, 2).AlignHorizontal(lipgloss.Left),
		Form:       lipgloss.NewStyle().Margin(1, 2),
		View:       lipgloss.NewStyle().Align(lipgloss.Center),
	}
}
