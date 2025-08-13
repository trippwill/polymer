package menu

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct {
	id string
}

func NewQuitAtom() QuitAtom {
	return QuitAtom{
		id: util.NewUniqeTypeId[QuitAtom](),
	}
}

var _ tea.Model = QuitAtom{}

func (q QuitAtom) Init() tea.Cmd                           { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                            { return "Goodbye!\n" }
