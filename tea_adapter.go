package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
)

// TeaAdapter adapts an Atom to the tea.Model interface.
type TeaAdapter struct {
	atom Atom
}

func WrapAtom(atom Atom) TeaAdapter {
	return TeaAdapter{
		atom: atom,
	}
}

// AdaptAtom wraps an Atom as a TeaAdapter and returns it with the command.
func AdaptAtom(atom Atom, cmd tea.Cmd) (TeaAdapter, tea.Cmd) {
	return TeaAdapter{
		atom: atom,
	}, cmd
}

var _ tea.Model = TeaAdapter{}

// Init calls the Atom's Init if it implements Initializer.
func (t TeaAdapter) Init() tea.Cmd {
	return OptionalInit(t.atom)
}

// Update delegates to the Atom's Update method.
func (t TeaAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return AdaptAtom(t.atom.Update(msg)) }

// View delegates to the Atom's View method.
func (t TeaAdapter) View() string { return t.atom.View() }
