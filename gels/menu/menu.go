package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// Item is a menu item that contains an Atom and a description.
type Item struct {
	poly.Atom
	description string
}

func NewMenuItem(atom poly.Atom, description string) Item {
	return Item{
		Atom:        atom,
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

var _ poly.Atom = Model{}

func (m Model) Name() string { return m.name }

func (m Model) Init() tea.Cmd {
	return tea.WindowSize()
}

func (m Model) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected, ok := m.list.SelectedItem().(Item); ok && selected.Atom != nil {
				return m, poly.Push(selected.Atom)
			}

		case "esc":
			return m, poly.Pop()
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return m.list.View()
}
