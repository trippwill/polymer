package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Atom is the fundamental unit of application state.
type Atom interface {
	Name() string
	Update(tea.Msg) (Atom, tea.Cmd)
	View() string
}

// Initializer allows an Atom to provide an initialization command.
type Initializer interface {
	Atom
	Init() tea.Cmd
}

// AtomDecorator wraps an Atom for extension.
type AtomDecorator struct {
	Atom
}

// OptionalInit checks if the Atom implements Initializer and calls its Init method if it does.
func OptionalInit(atom Atom) tea.Cmd {
	if initializer, ok := atom.(Initializer); ok {
		return initializer.Init()
	}
	return nil
}

// DecorateAtom wraps an Atom in an AtomDecorator.
func DecorateAtom(atom Atom) AtomDecorator {
	return AtomDecorator{Atom: atom}
}

// Wrap wraps an Atom in the AtomDecorator, allowing for further decoration or extension.
func (d *AtomDecorator) Wrap(atom Atom) *AtomDecorator {
	d.Atom = atom
	return d
}
