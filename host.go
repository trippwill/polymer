package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

type Host struct {
	name  string
	state tea.Model
}

func NewHost(name string, root tea.Model, options ...LensOption) tea.Model {
	if root == nil {
		panic("root state cannot be nil")
	}

	if len(options) > 0 {
		root = NewLens(root, options...)
	}

	host := &Host{
		name:  name,
		state: root,
	}

	return host
}

var _ tea.Model = Host{}

// Init implements [tea.Model].
func (h Host) Init() tea.Cmd {
	return tea.Sequence(
		trace.TraceInfo(">>>> Initializing host: "+h.name),
		tea.SetWindowTitle(h.name),
		h.state.Init(),
		tea.WindowSize(),
	)
}

// Update implements [tea.Model].
func (h Host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
func (h Host) View() string {
	if h.state == nil {
		return ""
	}

	return h.state.View()
}
