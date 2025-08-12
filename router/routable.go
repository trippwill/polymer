package router

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

// Routable can update and render itself.
// Most tea.Models and bubbles implement this by default.
type Routable[T any] interface {
	Update(msg tea.Msg) (T, tea.Cmd)
	View() string
}

// ErrUnknownSlot signals an unknown slot.
var ErrUnknownSlot error = errors.New("unknown slot")
