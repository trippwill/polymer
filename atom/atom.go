package atom

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model is the fundamental unit of application state.
type Model interface {
	Name() string
	Update(tea.Msg) (Model, tea.Cmd)
	View() string
}

// Initializer allows an Atom to provide an initialization command.
type Initializer interface {
	Model
	Init() tea.Cmd
}

// AtomDecorator wraps an Atom for extension.
type AtomDecorator struct {
	Model
}

// OptionalInit checks if the Atom implements Initializer and calls its Init method if it does.
func OptionalInit(atom Model) tea.Cmd {
	if initializer, ok := atom.(Initializer); ok {
		return initializer.Init()
	}
	return nil
}

// DecorateAtom wraps an Atom in an AtomDecorator.
func DecorateAtom(atom Model) AtomDecorator {
	return AtomDecorator{Model: atom}
}

// Wrap wraps an Atom in the AtomDecorator, allowing for further decoration or extension.
func (d *AtomDecorator) Wrap(atom Model) *AtomDecorator {
	d.Model = atom
	return d
}
