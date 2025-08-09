// Package menu provides a simple menu system for selecting and activating [poly.Atomic]s.
package menu

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Menu displays a list of [Item]s and activates the selected one.
type Menu[X any] struct {
	router.Router[list.Model, poly.Atomic[X]]
	id  string
	log trace.Tracer
}

// NewMenu creates a new [Menu] with the given title and items.
func NewMenu[T any](title string, items ...Item[T]) Menu[T] {
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
	return Menu[T]{
		Router: router.NewRouter[list.Model, poly.Atomic[T]](displayList, nil, router.SlotT),
		id:     id,
		log:    trace.NewTracerWithId(trace.CategoryMenu, id),
	}
}

// ConfigureList configures the underlying list model of the Menu.
func (m *Menu[T]) ConfigureList(fn func(*list.Model)) {
	if fn != nil {
		m.log.Trace("Configuring list model with custom function")
		m.ApplyT(fn)
	} else {
		m.log.Warn("No configuration function provided for list model")
	}
}

// Id implements [poly.Identifier].
func (m Menu[I]) Id() string { return m.id }

// Init implements [poly.Initializer].
func (m Menu[T]) Init() tea.Cmd { return tea.WindowSize() }

var _ poly.Atomic[any] = Menu[any]{}

// SetContext implements [poly.Atomic].
func (m Menu[I]) SetContext(context I) poly.Atomic[I] {
	// no-op for Menu, as it doesn't use context
	return m
}

// Update implements [poly.Atomic].
func (m Menu[T]) Update(msg tea.Msg) (poly.Atomic[T], tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ApplyT(func(slot *list.Model) {
			m.log.Trace("Updating window size for list model")
			slot.SetSize(msg.Width, msg.Height)
		})

	case tea.KeyMsg:
		if m.GetTarget() == router.SlotT {
			switch msg.String() {
			case "enter":
				if selected, ok := m.GetSlotAsT().SelectedItem().(Item[T]); ok && selected != nil {
					m.SlotU = selected
					m.SetTarget(router.SlotU)
					m.log.Trace("Selected item: %s", selected.Title())
					return m, poly.OptionalInit(selected)
				}

			case "esc":
				m.log.Trace("Escape key pressed, going back")
				return nil, nil
			}
		}
	}

	if m.GetSlotAsU() == nil {
		m.SetTarget(router.SlotT)
	}

	var cmd tea.Cmd
	m.Router, cmd = m.Route(msg)
	return m, cmd
}

// View implements [poly.Atomic].
func (m Menu[T]) View() string {
	return m.Render()
}
