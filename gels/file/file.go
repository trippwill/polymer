package file

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// FileType represents what can be selected
type FileType int

const (
	FilesOnly FileType = iota
	DirsOnly
	FilesAndDirs
)

// Config holds configuration for the file selector
type Config struct {
	Title       string
	FileType    FileType
	CurrentDir  string
	ShowHidden  bool
}

// Selector wraps bubbles/filepicker for Polymer integration
type Selector struct {
	filepicker filepicker.Model
	config     Config
	name       string
}

// NewSelector creates a new file selector wrapping bubbles/filepicker
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

var _ poly.Atom = (*Selector)(nil)

func (s *Selector) Name() string { 
	return s.name 
}

func (s *Selector) Init() tea.Cmd {
	return s.filepicker.Init()
}

func (s *Selector) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return s, poly.Pop()
		}
	}

	var cmd tea.Cmd
	s.filepicker, cmd = s.filepicker.Update(msg)

	// Check if user selected a file
	if didSelect, path := s.filepicker.DidSelectFile(msg); didSelect {
		var selectionType string
		switch s.config.FileType {
		case FilesOnly:
			selectionType = "file"
		case DirsOnly:
			selectionType = "directory"
		default:
			selectionType = "file"
		}
		
		return s, tea.Sequence(
			poly.FileSelection([]string{path}, selectionType),
			poly.Pop(),
		)
	}

	return s, cmd
}

func (s *Selector) View() string {
	return s.filepicker.View()
}