package file

import tea "github.com/charmbracelet/bubbletea"

// SelectionType represents the type of file selection
//
//go:generate stringer -type=SelectionType -trimprefix=SelectionType
type SelectionType int

const (
	SelectionTypeFile SelectionType = iota
	SelectionTypeDirectory
	SelectionTypeFiles
	SelectionTypeDirectories
	SelectionTypeMixed
)

// FileSelectionMsg represents a file selection result
type FileSelectionMsg struct {
	Files []string
	Type  SelectionType
}

// FileSelection creates a command to send file selection results
func FileSelection(files []string, selectionType SelectionType) tea.Cmd {
	return func() tea.Msg {
		return FileSelectionMsg{
			Files: files,
			Type:  selectionType,
		}
	}
}
