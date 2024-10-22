package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	ComponentWidth  int
	ComponentHeight int

	Header       lipgloss.Style
	Conversation lipgloss.Style
	Timestamp    lipgloss.Style
	BubbleWrap   func(textLen int) lipgloss.Style

	LeftBubble lipgloss.Style
	leftAlign  lipgloss.Style

	RightBubble lipgloss.Style
	rightAlign  lipgloss.Style
}

func DefaultChatStyleFunc(width, height int) *ChatStyles {
	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"
	
	fullWidth := lipgloss.NewStyle().Width(width)

	return &ChatStyles{
		ComponentWidth:  width,
		ComponentHeight: height,

		// TODO: Add cleaner truncation for the header content. Truncation should not affect the line break effect.
		Header:       lipgloss.NewStyle().MaxWidth(width-10).MaxHeight(2).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Padding(0, 4),
		Conversation: lipgloss.NewStyle().Height(height - (6)), // accounting for the header and input heights
		Timestamp:    fullWidth.AlignHorizontal(lipgloss.Center),
		BubbleWrap: func(textLen int) lipgloss.Style {
			bubbleMaxWidth := (width / 5) * 4
			return lipgloss.NewStyle().Width(min(textLen, bubbleMaxWidth))
		},

		LeftBubble: lipgloss.NewStyle().Border(leftBubbleBorder, true),
		leftAlign:  fullWidth.AlignHorizontal(lipgloss.Left),

		RightBubble: lipgloss.NewStyle().Border(rightBubbleBorder, true),
		rightAlign:  fullWidth.AlignHorizontal(lipgloss.Right),
	}
}
