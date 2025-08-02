package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
	"github.com/trippwill/polymer/trace"
)

// Item is a menu item that contains an Atom and a description.
type Item struct {
	atom.Model
	description string
}

func NewMenuItem(atom atom.Model, description string) Item {
	return Item{
		Model:       atom,
		description: description,
	}
}

var _ list.DefaultItem = Item{}

func (m Item) Title() string       { return m.Name() }
func (m Item) Description() string { return m.description }
func (m Item) FilterValue() string { return fmt.Sprintf("%s %s", m.Name(), m.description) }

// Model displays a list of options and activates the selected Atom.
type Model struct {
	list list.Model
	name string
}

// NewMenu creates a new Menu with the given title and items.
func NewMenu(title string, items ...Item) *Model {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select item"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "go back"),
			),
		}
	}

	return &Model{list: l, name: title}
}

var _ atom.Model = Model{}

func (m Model) Name() string { return m.name }

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected, ok := m.list.SelectedItem().(Item); ok && selected.Model != nil {
				return m, tea.Sequence(
					trace.TraceInfo("selected: "+selected.Name()),
					atom.Push(selected.Model),
				)
			}

		case "esc":
			return m, atom.Pop()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
