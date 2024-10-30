package views

import "github.com/charmbracelet/lipgloss"

type AppStyles struct {
	View          lipgloss.Style
	Conversations lipgloss.Style
	Chat          lipgloss.Style
}

func DefaultAppStyles() *AppStyles {
	return &AppStyles{
		View:          lipgloss.NewStyle().Margin(1).Align(lipgloss.Center, lipgloss.Center),
		Conversations: lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderRight(true).BorderForeground(lipgloss.Color("237")).PaddingRight(2).MarginRight(2),
		Chat:          lipgloss.NewStyle(),
	}
}

func (s *AppStyles) TotalFrameSize() (int, int) {
	frameStyles := []lipgloss.Style{
		s.View,
		s.Conversations,
		s.Chat,
	}

	var width, height int
	var frameWidth, frameHeight int
	for _, style := range frameStyles {
		frameWidth, frameHeight = style.GetFrameSize()
		width += frameWidth
		height += frameHeight
	}

	return width, height
}
