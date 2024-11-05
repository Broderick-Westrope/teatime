package components

import "github.com/charmbracelet/lipgloss"

var (
	// The border to use for a bubble aligned to the left of the messages.
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
	// The border to use for a bubble aligned to the right of the messages.
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

type chatStyles struct {
	Width  int
	Height int

	Header       lipgloss.Style
	Conversation lipgloss.Style
	Timestamp    lipgloss.Style

	BubbleStyleFunc func(value string, alignRight bool, textLen int) string

	InputPrompt      lipgloss.Style
	InputText        lipgloss.Style
	InputPlaceholder lipgloss.Style
	InputCompletion  lipgloss.Style
	InputCursor      lipgloss.Style
}

type chatStyleFunc func(width, height int) *chatStyles

func enabledChatStyleFunc(width, height int) *chatStyles {
	leftBubble := lipgloss.NewStyle().Padding(0, 1).Border(leftBubbleBorder, true).BorderForeground(lipgloss.Color("6"))
	rightBubble := lipgloss.NewStyle().Padding(0, 1).Border(rightBubbleBorder, true).BorderForeground(lipgloss.Color("5"))
	fullWidth := lipgloss.NewStyle().Width(width)
	leftAlign := fullWidth.AlignHorizontal(lipgloss.Left)
	rightAlign := fullWidth.AlignHorizontal(lipgloss.Right)
	bubbleMaxWidth := (width / 10) * 7

	return &chatStyles{
		Width:  width,
		Height: height,

		// this width is a hard limit, but truncation should be used better manage the header width
		Header: lipgloss.NewStyle().MaxWidth(width-10).MaxHeight(2).
			BorderStyle(lipgloss.NormalBorder()).BorderBottom(true).Padding(0, 4),
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

		InputPrompt:      lipgloss.NewStyle(),
		InputText:        lipgloss.NewStyle(),
		InputPlaceholder: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		InputCompletion:  lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		InputCursor:      lipgloss.NewStyle(),
	}
}

func disabledChatStyleFunc(width, height int) *chatStyles {
	styles := enabledChatStyleFunc(width, height)

	disabledForeground := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).BorderForeground(lipgloss.Color("240"))

	styles.Header = styles.Header.Inherit(disabledForeground)
	styles.Timestamp = styles.Timestamp.Inherit(disabledForeground)

	leftBubble := lipgloss.NewStyle().Padding(0, 1).Border(leftBubbleBorder, true).Inherit(disabledForeground)
	rightBubble := lipgloss.NewStyle().Padding(0, 1).Border(rightBubbleBorder, true).Inherit(disabledForeground)
	fullWidth := lipgloss.NewStyle().Width(width).Inherit(disabledForeground)
	leftAlign := fullWidth.AlignHorizontal(lipgloss.Left)
	rightAlign := fullWidth.AlignHorizontal(lipgloss.Right)
	bubbleMaxWidth := (width / 10) * 7

	styles.BubbleStyleFunc = func(value string, alignRight bool, textLen int) string {
		value = lipgloss.NewStyle().Width(min(textLen, bubbleMaxWidth)).Inherit(disabledForeground).Render(value)

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

	styles.InputPrompt = lipgloss.NewStyle().Inherit(disabledForeground)
	styles.InputText = lipgloss.NewStyle().Inherit(disabledForeground)
	styles.InputPlaceholder = lipgloss.NewStyle().Inherit(disabledForeground)
	styles.InputCompletion = lipgloss.NewStyle().Inherit(disabledForeground)
	styles.InputCursor = lipgloss.NewStyle().Inherit(disabledForeground)

	return styles
}
