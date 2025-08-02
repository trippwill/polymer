package polymer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

// AtomAdapter adapts a tea.Model to the Atom interface.
type AtomAdapter struct {
	Model    tea.Model
	AtomName string
}

// AdaptModel wraps a tea.Model as an Atom if needed and returns it with the command.
func AdaptModel(mode tea.Model, cmd tea.Cmd) (atom.Model, tea.Cmd) {
	return AtomAdapter{Model: mode}, cmd
}

var _ atom.Model = AtomAdapter{}

// Name returns the AtomAdapter's name or the type name of the model.
func (a AtomAdapter) Name() string {
	if a.AtomName != "" {
		return a.AtomName
	}

	return fmt.Sprintf("%T", a.Model)
}

// Init calls the underlying model's Init method.
func (a AtomAdapter) Init() tea.Cmd {
	return a.Model.Init()
}

// Update delegates to the model's Update and adapts the result.
func (a AtomAdapter) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	return AdaptModel(a.Model.Update(msg))
}

// View delegates to the model's View method.
func (a AtomAdapter) View() string { return a.Model.View() }
