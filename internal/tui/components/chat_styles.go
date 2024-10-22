package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	LeftBubble  lipgloss.Style
	RightBubble lipgloss.Style
	leftAlign   lipgloss.Style
	rightAlign  lipgloss.Style
}

func DefaultChatStyleFunc(width, _ int) *ChatStyles {
	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"

	bubbleDimensions := lipgloss.NewStyle().Width((width / 5) * 4)
	alignDimensions := lipgloss.NewStyle().Width(width)

	return &ChatStyles{
		LeftBubble:  bubbleDimensions.Border(leftBubbleBorder, true),
		RightBubble: bubbleDimensions.Border(rightBubbleBorder, true),
		leftAlign:   alignDimensions.AlignHorizontal(lipgloss.Left),
		rightAlign:  alignDimensions.AlignHorizontal(lipgloss.Right),
	}
}
