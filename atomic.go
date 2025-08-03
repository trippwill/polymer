package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// HasName provides the Name method.
type HasName interface {
	Name() string
}

type Named struct {
	name string
}

func (n Named) Name() string { return n.name }

// HasId provides the Id method.
type HasId interface {
	Id() uint32
}

type Identified struct {
	id uint32
}

func (i Identified) Id() uint32 { return i.id }

type Atomic interface {
	tea.Model
	HasName
	HasId
}

type teaAtomic struct {
	tea.Model
	Named
	Identified
}

func NewAtomicTea(model tea.Model, name string) Atomic {
	return &teaAtomic{
		Model:      model,
		Named:      Named{name: name},
		Identified: Identified{id: util.NewId()},
	}
}

type AtomicProxy struct {
	Named
	Identified
}

func NewAtomicProxy(name string) Atomic {
	return &AtomicProxy{
		Named:      Named{name: name},
		Identified: Identified{id: util.NewId()},
	}
}

func (AtomicProxy) Init() tea.Cmd                           { return nil }
func (AtomicProxy) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return nil, nil }
func (ap AtomicProxy) View() string                         { return ap.Name() + " (Proxy)\n" }
