package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// Atom is an embeddable struct providing
// a partial default implementation of [Atomic].
//
// Components must implement [tea.Model.Update] and [tea.Model.View].
// Optionally, they can also implement [tea.Model.Init].
type Atom struct {
	id   uint32
	name string
}

// NewAtom creates a new [Atom] with the given name.
func NewAtom(name string) Atom {
	return Atom{
		id:   util.NewId(),
		name: name,
	}
}

// OverrideID overrides the current ID of the atom.
func (atom *Atom) OverrideID(id uint32) *Atom {
	if id == 0 {
		panic("ID 0 is reserved and cannot be used")
	}

	atom.id = id
	return atom
}

func (atom Atom) Name() string { return atom.name }

func (atom Atom) Id() uint32 { return atom.id }

// Init implements [tea.Model].
// It is a no-op and returns nil.
func (atom Atom) Init() tea.Cmd { return nil }
