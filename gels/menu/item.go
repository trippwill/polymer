package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
)

// Item is an interface for items that can be displayed in a Menu.
type Item[T any] interface {
	tea.Model
	list.DefaultItem
	Description() string // Description of the item.
}

// ItemAdapter wraps a [tea.Model] providind a simple implementation of Item
// by adding title and description fields, for use in a [Menu].
type ItemAdapter[T any] struct {
	tea.Model
	title       string
	description string
}

// AdaptItem creates a new [ItemAdapter] with the given model, title, and description.
func AdaptItem[T any](model tea.Model, title, description string) *ItemAdapter[T] {
	return &ItemAdapter[T]{
		Model:       model,
		description: description,
		title:       title,
	}
}

var _ Item[any] = ItemAdapter[any]{}

// Implement [list.DefaultItem] for SimpleItem.
func (e ItemAdapter[T]) Title() string       { return e.title }
func (e ItemAdapter[T]) Description() string { return e.description }
func (e ItemAdapter[T]) FilterValue() string { return fmt.Sprintf("%s %s", e.title, e.description) }

// Update implements [tea.Model].
func (e ItemAdapter[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e.Model.Update(msg)
}

// View implements [tea.Model].
func (e ItemAdapter[T]) View() string {
	return e.Model.View()
}

// Init implements [tea.Model].
func (e ItemAdapter[T]) Init() tea.Cmd { return e.Model.Init() }

var _ atoms.ContextAware[any] = (*ItemAdapter[any])(nil)

// SetContext implements [atoms.ContextAware].
func (e *ItemAdapter[T]) SetContext(ctx T) {
	if contextAware, ok := e.Model.(atoms.ContextAware[T]); ok {
		contextAware.SetContext(ctx)
	}
}
