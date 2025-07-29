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

// FileSelectionMsg represents a file selection result
type FileSelectionMsg struct {
	Files []string
	Type  string // "file", "directory", or "files"
}

// FileSelection creates a command to send file selection results
func FileSelection(files []string, selectionType string) tea.Cmd {
	return func() tea.Msg {
		return FileSelectionMsg{
			Files: files,
			Type:  selectionType,
		}
	}
}
