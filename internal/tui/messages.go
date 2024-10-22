package tui

type ComponentSizeMsg struct {
	Width  int
	Height int
}

type FatalErrorMsg error
