package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

// TeaAdapter adapts an Atom to the tea.Model interface.
type TeaAdapter struct {
	atom.Model
}

func WrapAtom(atom atom.Model) TeaAdapter {
	return TeaAdapter{
		Model: atom,
	}
}

// AdaptAtom wraps an Atom as a TeaAdapter and returns it with the command.
func AdaptAtom(atom atom.Model, cmd tea.Cmd) (TeaAdapter, tea.Cmd) {
	return TeaAdapter{
		Model: atom,
	}, cmd
}

var _ tea.Model = TeaAdapter{}

// Init calls the Atom's Init if it implements Initializer.
func (t TeaAdapter) Init() tea.Cmd {
	return atom.OptionalInit(t.Model)
}

// Update delegates to the Atom's Update method.
func (t TeaAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return AdaptAtom(t.Model.Update(msg)) }

// View delegates to the Atom's View method.
func (t TeaAdapter) View() string { return t.Model.View() }
