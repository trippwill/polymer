package util

import (
	tea "github.com/charmbracelet/bubbletea"
)

var currentId uint32 = 1

// NewId generates a new unique identifier.
func NewId() uint32 {
	currentId++
	localId := currentId
	return localId
}

// Broadcast sends a message to the event loop.
func Broadcast[T any](msg T) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// ContextMsg carries context data.
type ContextMsg[T any] struct {
	Context T
}

// ContextUpdate sends a context message.
func ContextUpdate[T any](ctx T) tea.Cmd {
	return func() tea.Msg {
		return ContextMsg[T]{Context: ctx}
	}
}
