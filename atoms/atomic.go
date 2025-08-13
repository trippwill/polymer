package atoms

import tea "github.com/charmbracelet/bubbletea"

// NilInit provided a default implementation of the [tea.Model] Init method.
type NilInit struct{}

// WindowSizeInit provides a default implementation of the [tea.Model] Init method.
type WindowSizeInit struct{}

// ContextAware is an interface for components that can set a context.
type ContextAware[X any] interface {
	SetContext(context X)
}

// ContextMsg is a message type that carries a context value.
type ContextMsg[X any] struct {
	Context X
}

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
