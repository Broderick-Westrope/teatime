package views

import "github.com/charmbracelet/lipgloss"

type lockStyles struct {
	Title lipgloss.Style
	Logo  lipgloss.Style
	Form  lipgloss.Style
	View  lipgloss.Style
}

func defaultLockStyles() *lockStyles {
	return &lockStyles{
		Title: lipgloss.NewStyle().Margin(1, 2).AlignHorizontal(lipgloss.Left),
		Logo:  lipgloss.NewStyle().Margin(1, 2).AlignHorizontal(lipgloss.Left),
		Form:  lipgloss.NewStyle().Margin(1, 2),
		View:  lipgloss.NewStyle().Align(lipgloss.Center),
	}
}

func (s *lockStyles) GetFormRelativeHorizontalFrameSize() int {
	return s.Logo.GetHorizontalFrameSize() + s.Form.GetHorizontalFrameSize()
}

func (s *lockStyles) GetFormRelativeVerticalFrameSize() int {
	return s.Title.GetVerticalFrameSize() + s.Form.GetVerticalFrameSize()
}
