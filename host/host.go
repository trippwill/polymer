package host

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

type Host struct {
	name  string
	state tea.Model
	log   trace.Tracer
}

func New(name string, root tea.Model) tea.Model {
	if root == nil {
		panic("root state cannot be nil")
	}

	host := &Host{
		name:  name,
		state: root,
		log:   trace.NewTracer(trace.CategoryHost),
	}

	return host
}

var _ tea.Model = Host{}

// Init implements [tea.Model].
func (h Host) Init() tea.Cmd {
	h.log.Info(">>>> Initializing host: " + h.name)
	return tea.Sequence(
		tea.SetWindowTitle(h.name),
		h.state.Init(),
	)
}

// Update implements [tea.Model].
func (h Host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.log.Trace("Host received message: %T {%v}", msg, msg)
	switch msg := msg.(type) {
	case error:
		log.Fatal("Error in application:", msg)

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
func (h Host) View() string {
	if h.state == nil {
		h.log.Warn("Host %s state is nil, returning empty view.", h.name)
		return ""
	}

	return h.state.View()
}
