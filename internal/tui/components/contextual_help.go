package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type ContextualHelp[T comparable] struct {
	keymaps map[T]KeyMap
	Current T

	Width   int
	ShowAll bool // if true, render the "full" help menu

	ShortSeparator string
	FullSeparator  string

	// The symbol we use in the short help when help items have been truncated
	// due to width. Periods of ellipsis by default.
	Ellipsis string

	Styles *ContextualHelpStyles
}

type KeyMap interface {
	// ShortHelp returns a slice of bindings to be displayed in the short
	// version of the help. The help bubble will render help in the order in
	// which the help items are returned here.
	ShortHelp() []key.Binding

	// FullHelp returns an extended group of help items, grouped by columns.
	// The help bubble will render the help in the order in which the help
	// items are returned here.
	FullHelp() [][]key.Binding
}

func NewContextualHelp[T comparable](keymaps map[T]KeyMap, current T) *ContextualHelp[T] {
	return &ContextualHelp[T]{
		keymaps: keymaps,
		Current: current,

		ShortSeparator: " • ",
		FullSeparator:  "    ",
		Ellipsis:       "…",
		Styles:         defaultContextualHelpStyles(),
	}
}

func (m *ContextualHelp[T]) View() string {
	if m.ShowAll {
		return m.FullHelpView(m.keymaps[m.Current].FullHelp())
	}
	return m.ShortHelpView(m.keymaps[m.Current].ShortHelp())
}

// ShortHelpView renders a single line help view from a slice of keybindings.
// If the line is longer than the maximum width it will be gracefully
// truncated, showing only as many help items as possible.
func (m *ContextualHelp[T]) ShortHelpView(bindings []key.Binding) string {
	if len(bindings) == 0 {
		return ""
	}

	var b strings.Builder
	var totalWidth int
	separator := m.Styles.ShortSeparator.Inline(true).Render(m.ShortSeparator)

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if totalWidth > 0 && i < len(bindings) {
			sep = separator
		}

		str := sep +
			m.Styles.ShortKey.Inline(true).Render(kb.Help().Key) + " " +
			m.Styles.ShortDesc.Inline(true).Render(kb.Help().Desc)

		w := lipgloss.Width(str)

		// If adding this help item would go over the available width, stop
		// drawing.
		if m.Width > 0 && totalWidth+w > m.Width {
			// Although if there's room for an ellipsis, print that.
			tail := " " + m.Styles.Ellipsis.Inline(true).Render(m.Ellipsis)
			tailWidth := lipgloss.Width(tail)

			if totalWidth+tailWidth < m.Width {
				b.WriteString(tail)
			}

			break
		}

		totalWidth += w
		b.WriteString(str)
	}

	return b.String()
}

// FullHelpView renders help columns from a slice of key binding slices. Each
// top level slice entry renders into a column.
func (m *ContextualHelp[T]) FullHelpView(groups [][]key.Binding) string {
	if len(groups) == 0 {
		return ""
	}

	// Linter note: at this time we don't think it's worth the additional
	// code complexity involved in preallocating this slice.
	//nolint:prealloc
	var (
		out []string

		totalWidth int
		sep        = m.Styles.FullSeparator.Render(m.FullSeparator)
		sepWidth   = lipgloss.Width(sep)
	)

	// Iterate over groups to build columns
	for i, group := range groups {
		if group == nil || !m.shouldRenderColumn(group) {
			continue
		}

		var (
			keys         []string
			descriptions []string
		)

		// Separate keys and descriptions into different slices
		for _, kb := range group {
			if !kb.Enabled() {
				continue
			}
			keys = append(keys, kb.Help().Key)
			descriptions = append(descriptions, kb.Help().Desc)
		}

		col := lipgloss.JoinHorizontal(lipgloss.Top,
			m.Styles.FullKey.Render(strings.Join(keys, "\n")),
			m.Styles.FullKey.Render(" "),
			m.Styles.FullDesc.Render(strings.Join(descriptions, "\n")),
		)

		// Column
		totalWidth += lipgloss.Width(col)
		if m.Width > 0 && totalWidth > m.Width {
			break
		}

		out = append(out, col)

		// Separator
		if i < len(group)-1 {
			totalWidth += sepWidth
			if m.Width > 0 && totalWidth > m.Width {
				break
			}
			out = append(out, sep)
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, out...)
}

func (m *ContextualHelp[T]) shouldRenderColumn(b []key.Binding) (ok bool) {
	for _, v := range b {
		if v.Enabled() {
			return true
		}
	}
	return false
}
