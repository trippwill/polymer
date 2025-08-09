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

// Slot identifies a [Routable] in a [Routed] or [Routed3] struct.
//
//go:generate stringer -type=Slot
type Slot uint8

const (
	SlotSkip Slot = iota // Skip: no slot handles the message.
	SlotT
	SlotU
	SlotV
)

// ErrUnknownSlot signals an unknown slot.
var ErrUnknownSlot error = errors.New("unknown slot")
