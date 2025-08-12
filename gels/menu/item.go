package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
)

// Item is an interface for items that can be displayed in a Menu.
type Item[T any] interface {
	poly.Atomic[T] // Must be an Atomic[T] to be used in a Menu.
	list.DefaultItem
	Description() string // Description of the item.
}

// ItemAdapter wraps an Atomic[T] and provides a simple implementation of Item
// by adding title and description fields.
type ItemAdapter[T any] struct {
	poly.Atomic[T] // Must be an Atomic[T] to be used in a Menu.
	title          string
	description    string
}

// AdaptItem creates a new [ItemAdapter] with the given model, title, and description.
func AdaptItem[T any](model poly.Atomic[T], title, description string) *ItemAdapter[T] {
	return &ItemAdapter[T]{
		Atomic:      model,
		description: description,
		title:       title,
	}
}

var _ Item[any] = ItemAdapter[any]{}

// Implement [list.DefaultItem] for SimpleItem.
func (e ItemAdapter[T]) Title() string       { return e.title }
func (e ItemAdapter[T]) Description() string { return e.description }
func (e ItemAdapter[T]) FilterValue() string { return fmt.Sprintf("%s %s", e.title, e.description) }

// Update implements [poly.Atomic].
func (e ItemAdapter[T]) Update(msg tea.Msg) (poly.Atomic[T], tea.Cmd) {
	return e.Atomic.Update(msg)
}

// View implements [poly.Atomic].
func (e ItemAdapter[T]) View() string {
	return e.Atomic.View()
}

var _ poly.ContextAware[any] = (*ItemAdapter[any])(nil)

// SetContext implements [poly.ContextAware].
func (e *ItemAdapter[T]) SetContext(ctx T) {
	if contextAware, ok := e.Atomic.(poly.ContextAware[T]); ok {
		contextAware.SetContext(ctx)
	}
}

// Init implements [poly.Initializer].
func (e ItemAdapter[T]) Init() tea.Cmd {
	return poly.OptionalInit(e.Atomic)
}
