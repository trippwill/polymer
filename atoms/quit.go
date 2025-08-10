package atoms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/util"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct {
	id string
}

func NewQuitAtom() poly.Atomic[any] {
	return QuitAtom{
		id: util.NewUniqeTypeId[QuitAtom](),
	}
}

var _ poly.Atomic[any] = (*QuitAtom)(nil)

func (q QuitAtom) Init() tea.Cmd                                  { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                                   { return "Goodbye!\n" }
func (q QuitAtom) SetContext(ctx any) poly.Atomic[any] {
	// No context needed for QuitAtom
	return q
}
