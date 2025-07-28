package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Host struct {
	name  string
	state Atom
}

func NewHost(root Atom, name string) *Host {
	return &Host{
		name:  name,
		state: root,
	}
}

var _ tea.Model = Host{}

func (h Host) Init() tea.Cmd {
	return OptionalInit(h.state)
}

func (h Host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return h, tea.Quit
		}
	}

	next, cmd := h.state.Update(msg)
	h.state = next
	return h, cmd
}

func (h Host) View() string {
	return h.state.View()
}
