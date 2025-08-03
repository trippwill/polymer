package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
	"github.com/trippwill/polymer/gels/file"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/trace"
)

// SelectionHandler wraps a menu to handle file selection results
type SelectionHandler struct {
	poly.Atom
	state         tea.Model
	lastSelection string
}

func NewSelectionHandler(title string, items ...menu.Item) *SelectionHandler {
	return &SelectionHandler{
		Atom:  poly.NewAtom(title),
		state: menu.NewMenu(title, items...),
	}
}

var _ poly.Modal = (*SelectionHandler)(nil)

// GetCurrent implements [polymer.Modal].
func (s *SelectionHandler) GetCurrent() poly.Atomic {
	if s.lastSelection != "" {
		return s
	}

	if m, ok := s.state.(poly.Atomic); ok {
		return m
	}

	return nil
}

func (s *SelectionHandler) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case file.FileSelectionMsg:
		// Handle file selection results
		log.Default().Printf("File selection result: %v, %v", msg.Type, msg.Files)
		if len(msg.Files) == 1 {
			s.lastSelection = fmt.Sprintf("Selected %v: %s", msg.Type, msg.Files[0])
		} else {
			s.lastSelection = fmt.Sprintf("Selected %d %v: %s", len(msg.Files), msg.Type, strings.Join(msg.Files, ", "))
		}
		return s, trace.TraceInfo(s.lastSelection)
	}

	next, cmd := s.state.Update(msg)
	if next == nil {
		return nil, tea.Quit
	}

	s.lastSelection = "" // Clear last selection on new update
	s.state = next
	return s, cmd
}

func (s *SelectionHandler) View() string {
	view := s.state.View()
	if s.lastSelection != "" {
		view += "\n\n" + s.lastSelection
	}
	return view
}

type QuitAtom struct {
	poly.Atom
}

func (q QuitAtom) Init() tea.Cmd                           { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                            { return "Goodbye!\n" }

func main() {
	// Set up a standard logger for tracing
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	logger := log.New(f, "", log.Lmsgprefix)

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
	root := NewSelectionHandler(
		"File Selector Demo",
		menu.NewItem(singleFileSelector, "Single File Selection"),
		menu.NewItem(singleDirSelector, "Single Directory Selection"),
		menu.NewItem(multiFileSelector, "Multiple File Selection"),
		menu.NewItem(multiDirSelector, "Multiple Directory Selection"),
		menu.NewItem(mixedMultiSelector, "Mixed Selection (Files & Dirs)"),
		menu.NewItem(QuitAtom{Atom: poly.NewAtom("Quit")}, "Exit Application"),
	)

	// Create the host and start the Bubble Tea program
	host := poly.NewHost(
		"File Selector Example",
		root,
		poly.WithLifecycleLogging(logger)...,
	)

	p := tea.NewProgram(host)
	if _, err := p.Run(); err != nil {
		logger.Fatalf("Application failed: %v", err)
	}
}
