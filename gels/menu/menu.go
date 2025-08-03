// Package menu provides a simple menu system for selecting and activating Atoms.
package menu

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// Item is a menu item that contains an Atom and a description.
type Item struct {
	poly.Atomic
	description string
}

func NewItem(atom poly.Atomic, description string) Item {
	return Item{
		Atomic:      atom,
		description: description,
	}
}

var _ list.DefaultItem = Item{}

func (m Item) Title() string       { return m.Name() }
func (m Item) Description() string { return m.description }
func (m Item) FilterValue() string { return m.Name() + " " + m.description }

// Menu displays a list of options and activates the selected Atom.
type Menu struct {
	poly.Atom
	list     list.Model
	selected tea.Model
}

// NewMenu creates a new Menu with the given title and items.
func NewMenu(title string, items ...Item) *Menu {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = &item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title

	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "go back"),
			),
		}
	}

	return &Menu{
		Atom: poly.NewAtom(title),
		list: l,
	}
}

// ConfigureList configures the underlying list model of the Menu.
func (m *Menu) ConfigureList(fn func(*list.Model)) {
	if fn != nil {
		fn(&m.list)
	}
}

var _ poly.Modal = Menu{}

func (m Menu) GetCurrent() poly.Atomic {
	switch selected := m.selected.(type) {
	case nil:
		return m
	case poly.Atomic:
		return selected
	default:
		return poly.NewAtomicTea(selected, "Selected")
	}
}

var _ poly.Atomic = Menu{}

func (m Menu) Init() tea.Cmd { return tea.WindowSize() }

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		if m.selected == nil {
			switch msg.String() {
			case "enter":
				if selected, ok := m.list.SelectedItem().(*Item); ok && selected != nil {
					m.selected = selected.Atomic
					return m, m.selected.Init()
				}

			case "esc":
				return nil, nil
			}
		}
	}

	var cmd tea.Cmd
	if m.selected != nil {
		m.selected, cmd = m.selected.Update(msg)
		return m, cmd
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Menu) View() string {
	if m.selected != nil {
		return m.selected.View()
	}

	return m.list.View()
}
