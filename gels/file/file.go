package file

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

// FileType filters what types of files are selectable in the file picker
//
//go:generate stringer -type=FileType
type FileType int

const (
	FilesOnly FileType = iota
	DirsOnly
	FilesAndDirs
)

// Config holds configuration for the file selector
type Config struct {
	Title      string
	FileType   FileType
	CurrentDir string
	ShowHidden bool
}

// Selector is a file/directory selector.
type Selector struct {
	filepicker filepicker.Model
	config     Config
	name       string
}

// NewSelector creates a new file selector.
func NewSelector(config Config) *Selector {
	if config.CurrentDir == "" {
		if wd, err := os.Getwd(); err == nil {
			config.CurrentDir = wd
		} else {
			config.CurrentDir = "."
		}
	}

	if config.Title == "" {
		config.Title = "Select File"
	}

	fp := filepicker.New()

	// Configure filepicker based on Config
	fp.ShowHidden = config.ShowHidden
	fp.ShowSize = true
	fp.ShowPermissions = false
	fp.CurrentDirectory = config.CurrentDir

	// Set file/directory permissions based on FileType
	switch config.FileType {
	case FilesOnly:
		fp.FileAllowed = true
		fp.DirAllowed = false
	case DirsOnly:
		fp.FileAllowed = false
		fp.DirAllowed = true
	case FilesAndDirs:
		fp.FileAllowed = true
		fp.DirAllowed = true
	}

	return &Selector{
		filepicker: fp,
		config:     config,
		name:       config.Title,
	}
}

var _ atom.Model = Selector{}

func (s Selector) Name() string {
	return s.name
}

func (s Selector) Init() tea.Cmd {
	return tea.Sequence(
		s.filepicker.Init(),
		tea.WindowSize(),
	)
}

func (s Selector) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, atom.Pop()
		}
	case tea.WindowSizeMsg:
		s.filepicker.SetHeight(msg.Height - 2) // Leave space for title
	}

	var cmd tea.Cmd
	s.filepicker, cmd = s.filepicker.Update(msg)

	// Check if user selected a file
	if didSelect, path := s.filepicker.DidSelectFile(msg); didSelect {
		var selectionType SelectionType
		switch s.config.FileType {
		case FilesOnly:
			selectionType = SelectionTypeFile
		case DirsOnly:
			selectionType = SelectionTypeDirectory
		default:
			selectionType = SelectionTypeFile
		}

		return s, tea.Sequence(
			FileSelection([]string{path}, selectionType),
			atom.Pop(),
		)
	}

	return s, cmd
}

func (s Selector) View() string {
	return s.filepicker.View()
}
