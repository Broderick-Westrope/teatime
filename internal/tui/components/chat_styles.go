package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	Header          lipgloss.Style
	Conversation    lipgloss.Style
	Timestamp       lipgloss.Style
	BubbleStyleFunc func(value string, alignRight bool, textLen int) string
}

func DefaultChatStyleFunc(width, height int) *ChatStyles {
	bubbleMaxWidth := (width / 5) * 4
	fullWidth := lipgloss.NewStyle().Width(width)

	leftBubbleBorder := lipgloss.RoundedBorder()
	leftBubbleBorder.BottomLeft = "└"
	leftBubble := lipgloss.NewStyle().Border(leftBubbleBorder, true)
	leftAlign := fullWidth.AlignHorizontal(lipgloss.Left)

	rightBubbleBorder := lipgloss.RoundedBorder()
	rightBubbleBorder.BottomRight = "┘"
	rightBubble := lipgloss.NewStyle().Border(rightBubbleBorder, true)
	rightAlign := fullWidth.AlignHorizontal(lipgloss.Right)

	return &ChatStyles{
		// TODO: Add cleaner truncation for the header content. Truncation should not affect the line break effect.
		Header:       lipgloss.NewStyle().MaxWidth(width-10).MaxHeight(2).BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Padding(0, 4),
		Conversation: lipgloss.NewStyle().Height(height - (6)), // accounting for the header and input heights
		Timestamp:    fullWidth.AlignHorizontal(lipgloss.Center),

		BubbleStyleFunc: func(value string, alignRight bool, textLen int) string {
			value = lipgloss.NewStyle().Width(min(textLen, bubbleMaxWidth)).Render(value)

			switch alignRight {
			case true:
				value = rightBubble.Render(value)
				value = rightAlign.Render(value)

			case false:
				value = leftBubble.Render(value)
				value = leftAlign.Render(value)
			}
			return value
		},
	}
}
