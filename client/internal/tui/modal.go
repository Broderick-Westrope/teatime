package tui

import tea "github.com/charmbracelet/bubbletea"

type Modal interface {
	tea.Model
	SetSize(width, height int)
}
