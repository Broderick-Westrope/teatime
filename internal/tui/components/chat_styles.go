package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	LeftBubble  lipgloss.Style
	RightBubble lipgloss.Style
	leftAlign   lipgloss.Style
	rightAlign  lipgloss.Style
}

func DefaultChatStyles() *ChatStyles {
	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"

	return &ChatStyles{
		LeftBubble:  lipgloss.NewStyle().Border(leftBubbleBorder, true),
		RightBubble: lipgloss.NewStyle().Border(rightBubbleBorder, true),
		leftAlign:   lipgloss.NewStyle().AlignHorizontal(lipgloss.Left),
		rightAlign:  lipgloss.NewStyle().AlignHorizontal(lipgloss.Right),
	}
}
