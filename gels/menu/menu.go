// Package menu is a simple menu system for selecting and activating [atoms.Atomic]s.
package menu

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Menu displays a list of [Item]s and activates the selected one.
type Menu struct {
	atoms.WindowSizeInit
	router.Dual[list.Model, tea.Model]
	id  string
	log trace.Tracer
}

// NewMenu creates a new [Menu] with the given title and items.
func NewMenu(title string, items ...Item) *Menu {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	displayList := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	displayList.Title = title

	displayList.AdditionalShortHelpKeys = func() []key.Binding {
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

	id := util.NewUniqueId(title)
	return &Menu{
		Dual: router.NewDual[list.Model, tea.Model](displayList, nil, router.DualSlotA),
		id:   id,
		log:  trace.NewTracerWithId(trace.CategoryMenu, id),
	}
}

// ConfigureList configures the underlying list model of the Menu.
func (m *Menu) ConfigureList(fn func(*list.Model)) {
	if fn != nil {
		m.log.Trace("Configuring list model with custom function")
		m.Dual = m.ConfigureA(fn)
	} else {
		m.log.Warn("No configuration function provided for list model")
	}
}

// Id implements [atoms.Identifier].
func (m Menu) Id() string { return m.id }

var _ tea.Model = Menu{}

// Update implements [tea.Model].
func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Dual = m.ConfigureA(func(slot *list.Model) {
			m.log.Trace("Updating window size for list model")
			slot.SetSize(msg.Width, msg.Height)
		})

	case tea.KeyMsg:
		if m.Target() == router.DualSlotA {
			switch msg.String() {
			case "enter":
				if selected, ok := m.GetA().(list.Model).SelectedItem().(Item); ok {
					if selected == nil {
						m.log.Warn("Selected item is nil, cannot proceed")
						return m, nil
					}

					m.Dual = m.SetB(selected)
					m.SetTarget(router.DualSlotB)
					m.log.Trace("Selected item: %s", selected.Title())
					return m, selected.Init()
				}

			case "esc":
				m.log.Trace("Escape key pressed, going back")
				return nil, nil
			}
		}
	}

	if !m.IsSet(router.DualSlotB) {
		m.SetTarget(router.DualSlotA)
	}

	var cmd tea.Cmd
	m.Dual, cmd = m.Route(msg)
	return m, cmd
}

// View implements [tea.Model].
func (m Menu) View() string { return m.Render() }
