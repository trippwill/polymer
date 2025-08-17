package atoms

import tea "github.com/charmbracelet/bubbletea"

// NilInit provided a default implementation of the [tea.Model] Init method.
type NilInit struct{}

// WindowSizeInit provides a default implementation of the [tea.Model] Init method.
type WindowSizeInit struct{}

// Identifier is for components that have a unique identifier.
type Identifier interface {
	Id() string
}

// Init implements [tea.Model].
func (a NilInit) Init() tea.Cmd {
	return nil
}

// Init implements [tea.Model].
func (w WindowSizeInit) Init() tea.Cmd {
	return tea.WindowSize()
}
