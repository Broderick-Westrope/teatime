package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	Header    lipgloss.Style
	Timestamp lipgloss.Style

	LeftBubble lipgloss.Style
	leftAlign  lipgloss.Style

	RightBubble lipgloss.Style
	rightAlign  lipgloss.Style
}

func DefaultChatStyleFunc(width, _ int) *ChatStyles {
	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"

	bubbleWidth := lipgloss.NewStyle().Width((width / 5) * 4)
	fullWidth := lipgloss.NewStyle().Width(width)

	return &ChatStyles{
		// TODO: Add cleaner truncation for the header content. Truncation should not affect the line break effect.
		Header:    lipgloss.NewStyle().MaxWidth(width-10).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Padding(0, 4),
		Timestamp: fullWidth.AlignHorizontal(lipgloss.Center),

		LeftBubble: bubbleWidth.Border(leftBubbleBorder, true),
		leftAlign:  fullWidth.AlignHorizontal(lipgloss.Left),

		RightBubble: bubbleWidth.Border(rightBubbleBorder, true),
		rightAlign:  fullWidth.AlignHorizontal(lipgloss.Right),
	}
}
