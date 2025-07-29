package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Host struct {
	name  string
	state Atom
}

func NewHost(name string, root Atom, options ...LensOption) *Host {
	if root == nil {
		panic("root state cannot be nil")
	}

	if len(options) > 0 {
		root = LensWrap(root, options...)
	}

	return &Host{
		name:  name,
		state: root,
	}
}

var _ tea.Model = Host{}

func (h Host) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle(h.name),
		OptionalInit(h.state),
	)
}

func (h Host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	var cmd tea.Cmd
	h.state, cmd = h.state.Update(msg)
	return h, cmd
}

func (h Host) View() string {
	return h.state.View()
}
