package tui

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func UpdateTypedModel[T tea.Model](model *T, msg tea.Msg) (tea.Cmd, error) {
	var ok bool
	newModel, cmd := (*model).Update(msg)
	*model, ok = newModel.(T)
	if !ok {
		return nil, fmt.Errorf("failed to update model of type %q: %w", reflect.TypeOf(model), ErrInvalidTypeAssertion)
	}
	return cmd, nil
}

// Window Overlay (CREDIT: https://gist.github.com/ras0q/9bf5d81544b22302393f61206892e2cd) ----------------------------------------

// OverlayCenter writes the overlay string onto the background string such that the middle of the
// overlay string will be at the middle of the overlay will be at the middle of the background.
func OverlayCenter(bg string, overlay string) string {
	row := lipgloss.Height(bg) / 2
	row -= lipgloss.Height(overlay) / 2
	col := lipgloss.Width(bg) / 2
	col -= lipgloss.Width(overlay) / 2
	return Overlay(bg, overlay, row, col)
}

// Overlay writes the overlay string onto the background string at the specified row and column.
// In this case, the row and column are zero indexed.
func Overlay(bg string, overlay string, row int, col int) string {
	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, overlayLine := range overlayLines {
		targetRow := i + row

		bgLine := bgLines[targetRow]
		bgLineWidth := ansi.StringWidth(bgLine)

		if bgLineWidth < col {
			bgLine += strings.Repeat(" ", col-bgLineWidth) // Add padding
		}

		bgLeft := ansi.Truncate(bgLine, col, "")
		bgRight := truncateLeft(bgLine, col+ansi.StringWidth(overlayLine))

		bgLines[targetRow] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(bgLines, "\n")
}

// truncateLeft removes characters from the beginning of a line, considering ANSI escape codes.
func truncateLeft(line string, padding int) string {
	if strings.Contains(line, "\n") {
		panic("line must not contain newline")
	}

	wrapped := strings.Split(ansi.Hardwrap(line, padding, true), "\n")
	if len(wrapped) == 1 {
		return ""
	}

	var ansiStyle string
	// Regular expression to match ANSI escape codes.
	ansiStyles := regexp.MustCompile(`\x1b[[\d;]*m`).FindAllString(wrapped[0], -1)
	if len(ansiStyles) > 0 {
		ansiStyle = ansiStyles[len(ansiStyles)-1]
	}

	return ansiStyle + strings.Join(wrapped[1:], "")
}
