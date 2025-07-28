package polymer

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg error

func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

type ContextMsg[T any] struct {
	Context T
}

func ContextUpdate[T any](ctx T) tea.Cmd {
	return func() tea.Msg {
		return ContextMsg[T]{Context: ctx}
	}
}
