package util

import (
	"fmt"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

var current uint32 = 1

// GetCurrentId returns the current unique identifier without incrementing it.
func GetCurrentId() uint32 { return current }

// GetNextId returns the next unique identifier and increments the current ID.
func GetNextId() uint32 {
	local := current
	current++
	return local
}

// NewUniqueId generates a new unique identifier
// from the provided prefix and a unique number.
func NewUniqueId(prefix string) string {
	return fmt.Sprintf("%s#%d", prefix, GetNextId())
}

// NewUniqeTypeId generates a unique identifier based on the type of T and a unique number.
// This is useful for creating unique IDs for components or models.
func NewUniqeTypeId[T any]() string {
	return NewUniqueId(reflect.TypeOf((*T)(nil)).Elem().Name())
}

// Broadcast sends a message to the event loop.
func Broadcast[T any](msg T) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
