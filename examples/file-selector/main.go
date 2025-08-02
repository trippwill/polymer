package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
	"github.com/trippwill/polymer/atom"
	"github.com/trippwill/polymer/gels/file"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/trace"
)

// SelectionHandler wraps a menu to handle file selection results
type SelectionHandler struct {
	*menu.Model
	lastSelection string
}

func NewSelectionHandler(title string, items ...menu.Item) *SelectionHandler {
	return &SelectionHandler{
		Model: menu.NewMenu(title, items...),
	}
}

func (s *SelectionHandler) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case file.FileSelectionMsg:
		// Handle file selection results
		if len(msg.Files) == 1 {
			s.lastSelection = fmt.Sprintf("Selected %v: %s", msg.Type, msg.Files[0])
		} else {
			s.lastSelection = fmt.Sprintf("Selected %d %v: %s", len(msg.Files), msg.Type, strings.Join(msg.Files, ", "))
		}
		return s, nil
	}

	return s.Model.Update(msg)
}

func (s *SelectionHandler) View() string {
	view := s.Model.View()
	if s.lastSelection != "" {
		view += "\n\n" + s.lastSelection
	}
	return view
}

func (s *SelectionHandler) Name() string {
	return s.Model.Name()
}

type QuitAtom struct{}

func (q QuitAtom) Init() tea.Cmd                            { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (atom.Model, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                             { return "Goodbye!\n" }
func (q QuitAtom) Name() string                             { return "Quit" }

func main() {
	// Set up a standard logger for tracing
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	logger := log.New(f, "file-selector: ", log.LstdFlags)

	// Create file selectors with different configurations
	singleFileSelector := file.NewSelector(file.Config{
		Title:      "Select Single File",
		FileType:   file.FilesOnly,
		ShowHidden: false,
	})

	singleDirSelector := file.NewSelector(file.Config{
		Title:      "Select Directory",
		FileType:   file.DirsOnly,
		ShowHidden: false,
	})

	multiFileSelector := file.NewMultiSelector(file.Config{
		Title:      "Select Multiple Files",
		FileType:   file.FilesOnly,
		ShowHidden: false,
	})

	multiDirSelector := file.NewMultiSelector(file.Config{
		Title:      "Select Multiple Directories",
		FileType:   file.DirsOnly,
		ShowHidden: false,
	})

	mixedMultiSelector := file.NewMultiSelector(file.Config{
		Title:      "Select Files and Directories",
		FileType:   file.FilesAndDirs,
		ShowHidden: true, // Show hidden files for this demo
	})

	// Create the main menu
	root := atom.NewChain(
		NewSelectionHandler(
			"File Selector Demo",
			menu.NewMenuItem(singleFileSelector, "Single File Selection"),
			menu.NewMenuItem(singleDirSelector, "Single Directory Selection"),
			menu.NewMenuItem(multiFileSelector, "Multiple File Selection"),
			menu.NewMenuItem(multiDirSelector, "Multiple Directory Selection"),
			menu.NewMenuItem(mixedMultiSelector, "Mixed Selection (Files & Dirs)"),
			menu.NewMenuItem(QuitAtom{}, "Exit Application"),
		),
	)

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"File Selector Example",
		root,
		trace.WithLifecycleLogging(logger, trace.LevelTrace)...,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Application failed: %v", err)
	}
}
