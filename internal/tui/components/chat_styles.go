package components

import "github.com/charmbracelet/lipgloss"

type ChatStyles struct {
	Width  int
	Height int

	Header          lipgloss.Style
	Conversation    lipgloss.Style
	Timestamp       lipgloss.Style
	BubbleStyleFunc func(value string, alignRight bool, textLen int) string

	leftBubbleBorder  lipgloss.Border
	rightBubbleBorder lipgloss.Border
}

type ChatStyleFunc func(width, height int) *ChatStyles

var (
	leftBubbleBorder = lipgloss.Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "└",
		BottomRight:  "╯",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}
	rightBubbleBorder = lipgloss.Border{
		Top:          "─",
		Bottom:       "─",
		Left:         "│",
		Right:        "│",
		TopLeft:      "╭",
		TopRight:     "╮",
		BottomLeft:   "╰",
		BottomRight:  "┘",
		MiddleLeft:   "├",
		MiddleRight:  "┤",
		Middle:       "┼",
		MiddleTop:    "┬",
		MiddleBottom: "┴",
	}
)

func EnabledChatStyleFunc(width, height int) *ChatStyles {
	leftBubble := lipgloss.NewStyle().Border(leftBubbleBorder, true)
	rightBubble := lipgloss.NewStyle().Border(rightBubbleBorder, true)
	fullWidth := lipgloss.NewStyle().Width(width)
	leftAlign := fullWidth.AlignHorizontal(lipgloss.Left)
	rightAlign := fullWidth.AlignHorizontal(lipgloss.Right)
	bubbleMaxWidth := (width / 5) * 4

	return &ChatStyles{
		Width:  width,
		Height: height,

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

func DisabledStyleFunc(width, height int) *ChatStyles {
	styles := EnabledChatStyleFunc(width, height)

	disabledForegroundColor := lipgloss.Color("240")

	styles.Header = styles.Header.Foreground(disabledForegroundColor).BorderForeground(disabledForegroundColor)
	styles.Timestamp = styles.Timestamp.Foreground(disabledForegroundColor)

	leftBubble := lipgloss.NewStyle().Border(leftBubbleBorder, true).BorderForeground(disabledForegroundColor)
	rightBubble := lipgloss.NewStyle().Border(rightBubbleBorder, true).BorderForeground(disabledForegroundColor)
	fullWidth := lipgloss.NewStyle().Width(width).Foreground(disabledForegroundColor)
	leftAlign := fullWidth.AlignHorizontal(lipgloss.Left)
	rightAlign := fullWidth.AlignHorizontal(lipgloss.Right)
	bubbleMaxWidth := (width / 5) * 4

	styles.BubbleStyleFunc = func(value string, alignRight bool, textLen int) string {
		value = lipgloss.NewStyle().Width(min(textLen, bubbleMaxWidth)).Foreground(disabledForegroundColor).Render(value)

		switch alignRight {
		case true:
			value = rightBubble.Render(value)
			value = rightAlign.Render(value)

		case false:
			value = leftBubble.Render(value)
			value = leftAlign.Render(value)
		}
		return value
	}

	return styles
}
