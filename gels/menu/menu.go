// Package menu is a simple menu system for selecting and activating [atoms.Atomic]s.
package menu

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/router/multi"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Menu displays a list of [Item]s and activates the selected one.
type Menu[X any] struct {
	atoms.WindowSizeInit
	multi.Router[list.Model, tea.Model]
	ctx X
	id  string
	log trace.Tracer
}

// NewMenu creates a new [Menu] with the given title and items.
func NewMenu[X any](title string, items ...Item[X]) *Menu[X] {
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
	return &Menu[X]{
		Router: multi.NewRouter[list.Model, tea.Model](displayList, nil, multi.SlotT),
		id:     id,
		log:    trace.NewTracerWithId(trace.CategoryMenu, id),
	}
}

// ConfigureList configures the underlying list model of the Menu.
func (m *Menu[X]) ConfigureList(fn func(*list.Model)) {
	if fn != nil {
		m.log.Trace("Configuring list model with custom function")
		m.Router = m.ConfigureT(fn)
	} else {
		m.log.Warn("No configuration function provided for list model")
	}
}

// Id implements [atoms.Identifier].
func (m Menu[X]) Id() string { return m.id }

// SetContext implements [atoms.ContextAware].
func (m *Menu[X]) SetContext(context X) {
	m.ctx = context
}

var _ tea.Model = Menu[any]{}

// Update implements [tea.Model].
func (m Menu[X]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Router = m.ConfigureT(func(slot *list.Model) {
			m.log.Trace("Updating window size for list model")
			slot.SetSize(msg.Width, msg.Height)
		})

	case tea.KeyMsg:
		if m.Target() == multi.SlotT {
			switch msg.String() {
			case "enter":
				if selected, ok := m.GetT().(list.Model).SelectedItem().(Item[X]); ok {
					if selected == nil {
						m.log.Warn("Selected item is nil, cannot proceed")
						return m, nil
					}

					if contextAware, ok := selected.(atoms.ContextAware[X]); ok {
						m.log.Trace("Setting context for selected item: %s", selected.Title())
						contextAware.SetContext(m.ctx)
					}

					m.Router = m.SetU(selected)
					m.SetTarget(multi.SlotU)
					m.log.Trace("Selected item: %s", selected.Title())
					return m, selected.Init()
				}

			case "esc":
				m.log.Trace("Escape key pressed, going back")
				return nil, nil
			}
		}
	}

	if !m.IsSet(multi.SlotU) {
		m.SetTarget(multi.SlotT)
	}

	var cmd tea.Cmd
	m.Router, cmd = m.Route(msg)
	return m, cmd
}

// View implements [tea.Model].
func (m Menu[X]) View() string { return m.Render() }
