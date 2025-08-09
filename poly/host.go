package poly

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

type Host[T any] struct {
	name  string
	state Atomic[T]
	log   trace.Tracer
}

func NewHost[T any](name string, root Atomic[T]) tea.Model {
	if root == nil {
		panic("root state cannot be nil")
	}

	host := &Host[T]{
		name:  name,
		state: root,
		log:   trace.NewTracer(trace.CategoryHost),
	}

	return host
}

var _ tea.Model = Host[any]{}

// Init implements [tea.Model].
func (h Host[T]) Init() tea.Cmd {
	h.log.Info(">>>> Initializing host: " + h.name)
	return tea.Sequence(
		tea.SetWindowTitle(h.name),
		tea.WindowSize(),
		OptionalInit(h.state),
	)
}

// Update implements [tea.Model].
func (h Host[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.log.Trace("Host %s received message: %T", h.name, msg)
	switch msg := msg.(type) {
	case error:
		log.Fatal("Error in host:", msg)
	case *ContextMsg[T]:
		if h.state != nil {
			h.log.Debug("Host %s received context message: %v", h.name, msg.Context)
			h.state = h.state.SetContext(msg.Context)
		} else {
			h.log.Warn("Host %s received context message but state is nil, ignoring.", h.name)
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return h, tea.Quit
		}
	}

	var cmd tea.Cmd
	h.state, cmd = h.state.Update(msg)
	if h.state == nil {
		h.log.Debug("Host %s state is nil, quitting.", h.name)
		return nil, tea.Quit
	}

	return h, cmd
}

// View implements [tea.Model].
func (h Host[T]) View() string {
	if h.state == nil {
		return ""
	}

	return h.state.View()
}
