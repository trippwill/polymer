package file

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

// FileItem represents a file or directory in the list
type FileItem struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
}

func (f FileItem) Title() string { 
	if f.IsDir {
		return fmt.Sprintf("📁 %s", f.Name)
	}
	return fmt.Sprintf("📄 %s", f.Name)
}

func (f FileItem) Description() string { 
	if f.IsDir {
		return "Directory"
	}
	if f.Size > 0 {
		return fmt.Sprintf("File (%d bytes)", f.Size)
	}
	return "File"
}

func (f FileItem) FilterValue() string { 
	return f.Name 
}

var _ list.DefaultItem = FileItem{}

// Selector is a file/directory selector gel
type Selector struct {
	list       list.Model
	config     Config
	currentDir string
	name       string
	width      int
	height     int
}

// NewSelector creates a new file selector
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

	s := &Selector{
		config:     config,
		currentDir: config.CurrentDir,
		name:       config.Title,
	}

	s.setupList()
	return s
}

func (s *Selector) setupList() {
	items := s.loadDirectory()
	
	s.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	s.list.Title = fmt.Sprintf("%s - %s", s.config.Title, s.currentDir)
	
	s.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "select/enter"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "go back"),
			),
			key.NewBinding(
				key.WithKeys("backspace"),
				key.WithHelp("backspace", "parent dir"),
			),
		}
	}
}

func (s *Selector) loadDirectory() []list.Item {
	entries, err := os.ReadDir(s.currentDir)
	if err != nil {
		return []list.Item{}
	}

	var items []list.Item

	// Add parent directory entry if not at root
	if s.currentDir != "/" && s.currentDir != "." {
		items = append(items, FileItem{
			Name:  "..",
			Path:  filepath.Dir(s.currentDir),
			IsDir: true,
		})
	}

	var files []FileItem
	for _, entry := range entries {
		// Skip hidden files unless configured to show them
		if !s.config.ShowHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		isDir := entry.IsDir()
		
		// Filter based on file type configuration
		switch s.config.FileType {
		case FilesOnly:
			if isDir {
				continue
			}
		case DirsOnly:
			if !isDir {
				continue
			}
		case FilesAndDirs:
			// Include both
		}

		var size int64
		if !isDir {
			if info, err := entry.Info(); err == nil {
				size = info.Size()
			}
		}

		files = append(files, FileItem{
			Name:  entry.Name(),
			Path:  filepath.Join(s.currentDir, entry.Name()),
			IsDir: isDir,
			Size:  size,
		})
	}

	// Sort: directories first, then files, alphabetically
	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	for _, file := range files {
		items = append(items, file)
	}

	return items
}

var _ poly.Atom = (*Selector)(nil)

func (s *Selector) Name() string { 
	return s.name 
}

func (s *Selector) Init() tea.Cmd {
	return tea.WindowSize()
}

func (s *Selector) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetSize(msg.Width, msg.Height)
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if selected, ok := s.list.SelectedItem().(FileItem); ok {
				if selected.IsDir {
					// Navigate into directory
					s.currentDir = selected.Path
					s.setupList()
					if s.width > 0 && s.height > 0 {
						s.list.SetSize(s.width, s.height)
					}
					return s, nil
				} else {
					// File selected - return result via command or notification
					return s, tea.Sequence(
						poly.Notify(poly.InfoLevel, fmt.Sprintf("Selected file: %s", selected.Path)),
						poly.Pop(), // Go back after selection
					)
				}
			}

		case "backspace":
			// Navigate to parent directory
			parent := filepath.Dir(s.currentDir)
			if parent != s.currentDir { // Avoid infinite loop at root
				s.currentDir = parent
				s.setupList()
				if s.width > 0 && s.height > 0 {
					s.list.SetSize(s.width, s.height)
				}
				return s, nil
			}

		case "esc":
			return s, poly.Pop()
		}
	}

	var cmd tea.Cmd
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

func (s *Selector) View() string {
	return s.list.View()
}