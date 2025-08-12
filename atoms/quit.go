package atoms

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/util"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom[X any] struct {
	id string
}

func NewQuitAtom[X any]() poly.Atomic[X] {
	return QuitAtom[X]{
		id: util.NewUniqeTypeId[QuitAtom[X]](),
	}
}

var _ poly.Atomic[any] = QuitAtom[any]{}

func (q QuitAtom[X]) Init() tea.Cmd                                { return tea.Quit }
func (q QuitAtom[X]) Update(msg tea.Msg) (poly.Atomic[X], tea.Cmd) { return q, nil }
func (q QuitAtom[X]) View() string                                 { return "Goodbye!\n" }
