package util

import tea "github.com/charmbracelet/bubbletea"

// Broadcast sends a message to the event loop.
func Broadcast(msg any) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

// ContextMsg is a message that carries context data.
type ContextMsg[T any] struct {
	Context T
}

// ContextUpdate sends a context message.
func ContextUpdate[T any](ctx T) tea.Cmd {
	return func() tea.Msg {
		return ContextMsg[T]{Context: ctx}
	}
}
