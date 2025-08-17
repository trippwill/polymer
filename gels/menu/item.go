package menu

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Item is an interface for items that can be displayed in a Menu.
type Item interface {
	tea.Model
	list.DefaultItem
	Description() string // Description of the item.
}

// ItemAdapter wraps a [tea.Model] providind a simple implementation of Item
// by adding title and description fields, for use in a [Menu].
type ItemAdapter struct {
	tea.Model
	title       string
	description string
}

// AdaptItem creates a new [ItemAdapter] with the given model, title, and description.
func AdaptItem(model tea.Model, title, description string) *ItemAdapter {
	return &ItemAdapter{
		Model:       model,
		description: description,
		title:       title,
	}
}

var _ Item = ItemAdapter{}

// Implement [list.DefaultItem] for SimpleItem.
func (e ItemAdapter) Title() string       { return e.title }
func (e ItemAdapter) Description() string { return e.description }
func (e ItemAdapter) FilterValue() string { return fmt.Sprintf("%s %s", e.title, e.description) }

// Update implements [tea.Model].
func (e ItemAdapter) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return e.Model.Update(msg)
}

// View implements [tea.Model].
func (e ItemAdapter) View() string {
	return e.Model.View()
}

// Init implements [tea.Model].
func (e ItemAdapter) Init() tea.Cmd { return e.Model.Init() }
