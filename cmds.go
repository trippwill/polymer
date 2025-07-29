package polymer

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg error

type NotificationMsg struct {
	string
	level NotificationLevel
}

func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

func Notify(level NotificationLevel, msg string) tea.Cmd {
	return func() tea.Msg {
		return NotificationMsg{level: level, string: msg}
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
