package poly

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

type Host[T any] struct {
	name  string
	state Atomic[T]
}

func NewHost[T any](name string, root Atomic[T]) tea.Model {
	if root == nil {
		panic("root state cannot be nil")
	}

	host := &Host[T]{
		name:  name,
		state: root,
	}

	return host
}

var _ tea.Model = Host[any]{}

// Init implements [tea.Model].
func (h Host[T]) Init() tea.Cmd {
	return tea.Sequence(
		trace.TraceInfo(">>>> Initializing host: "+h.name),
		tea.SetWindowTitle(h.name),
		tea.WindowSize(),
		OptionalInit(h.state),
	)
}

// Update implements [tea.Model].
func (h Host[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("Host %s received message: %T\n", h.name, msg)
	switch msg := msg.(type) {
	case trace.TraceMsg:
		switch msg.Level {
		case trace.LevelTrace:
			log.Printf("TRACE: %s\n", msg.Msg)
		case trace.LevelDebug:
			log.Printf("DEBUG: %s\n", msg.Msg)
		case trace.LevelInfo:
			log.Printf("INFO: %s\n", msg.Msg)
		case trace.LevelWarn:
			log.Printf("WARN: %s\n", msg.Msg)
		}
	case error:
		log.Fatal("Error in host:", msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return h, tea.Quit
		}
	}

	var cmd tea.Cmd
	h.state, cmd = h.state.Update(msg)
	if h.state == nil {
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
