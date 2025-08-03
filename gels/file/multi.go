package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// SelectedFileItem represents a selected file in the selection list
type SelectedFileItem struct {
	Name string
	Path string
}

func (s SelectedFileItem) Title() string       { return s.Name }
func (s SelectedFileItem) Description() string { return s.Path }
func (s SelectedFileItem) FilterValue() string { return s.Name }

var _ list.DefaultItem = SelectedFileItem{}

// MultiSelector combines filepicker and list for multi-selection
type MultiSelector struct {
	poly.Atom
	filepicker       filepicker.Model
	selectedList     list.Model
	config           Config
	selected         map[string]SelectedFileItem // map of path -> item for selected items
	name             string
	showingSelection bool // toggle between filepicker and selection list
}

// NewMultiSelector creates a new multi-file selector
func NewMultiSelector(config Config) *MultiSelector {
	if config.Title == "" {
		config.Title = "Select Files"
	}

	// Set up filepicker
	fp := filepicker.New()
	fp.ShowHidden = config.ShowHidden
	fp.ShowSize = true
	fp.ShowPermissions = false

	if config.CurrentDir != "" {
		fp.CurrentDirectory = config.CurrentDir
	}

	// Configure file/directory permissions based on FileType
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

	// Set up selection list
	selectedList := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	selectedList.Title = "Selected Files"
	selectedList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("delete"),
				key.WithHelp("del", "remove"),
			),
			key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "toggle view"),
			),
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "confirm selection"),
			),
			key.NewBinding(
				key.WithKeys("esc"),
				key.WithHelp("esc", "cancel"),
			),
		}
	}

	return &MultiSelector{
		Atom:             poly.NewAtom(config.Title),
		filepicker:       fp,
		selectedList:     selectedList,
		config:           config,
		selected:         make(map[string]SelectedFileItem),
		showingSelection: false,
	}
}

func (ms *MultiSelector) updateSelectedList() {
	items := make([]list.Item, 0, len(ms.selected))
	for _, item := range ms.selected {
		items = append(items, item)
	}
	ms.selectedList.SetItems(items)
	ms.selectedList.Title = fmt.Sprintf("Selected Files (%d)", len(ms.selected))
}

func (ms *MultiSelector) addSelection(path, name string) {
	ms.selected[path] = SelectedFileItem{
		Name: name,
		Path: path,
	}
	ms.updateSelectedList()
}

func (ms *MultiSelector) removeSelection(path string) {
	delete(ms.selected, path)
	ms.updateSelectedList()
}

func (ms MultiSelector) getSelectedPaths() []string {
	var paths []string
	for path := range ms.selected {
		paths = append(paths, path)
	}
	return paths
}

func (ms MultiSelector) getSelectionType() SelectionType {
	if len(ms.selected) == 0 {
		return SelectionTypeFiles
	}

	// Determine selection type based on config and what was actually selected
	switch ms.config.FileType {
	case FilesOnly:
		if len(ms.selected) == 1 {
			return SelectionTypeFile
		}
		return SelectionTypeFiles
	case DirsOnly:
		if len(ms.selected) == 1 {
			return SelectionTypeDirectory
		}
		return SelectionTypeDirectories
	case FilesAndDirs:
		// For mixed mode, we need to check what types were actually selected
		hasFiles := false
		hasDirs := false

		for _, item := range ms.selected {
			if info, err := os.Stat(item.Path); err == nil {
				if info.IsDir() {
					hasDirs = true
				} else {
					hasFiles = true
				}
			}
		}

		if hasFiles && hasDirs {
			return SelectionTypeMixed
		} else if hasDirs {
			if len(ms.selected) == 1 {
				return SelectionTypeDirectory
			}
			return SelectionTypeDirectories
		} else {
			if len(ms.selected) == 1 {
				return SelectionTypeFile
			}
			return SelectionTypeFiles
		}
	default:
		return SelectionTypeFiles
	}
}

var _ poly.Atomic = MultiSelector{}

func (ms MultiSelector) Init() tea.Cmd {
	return tea.Sequence(
		ms.filepicker.Init(),
		tea.WindowSize(),
	)
}

func (ms MultiSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ms.filepicker.SetHeight(msg.Height)
		ms.selectedList.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			// Toggle between filepicker and selection views
			ms.showingSelection = !ms.showingSelection
			return ms, nil

		case "esc":
			if ms.showingSelection {
				ms.showingSelection = false
				return ms, nil
			}
			// Complete selection or go back
			if len(ms.selected) > 0 {
				paths := ms.getSelectedPaths()
				selectionType := ms.getSelectionType()
				return nil, FileSelection(paths, selectionType)
			}
			return ms, nil

		case "enter":
			if ms.showingSelection && len(ms.selected) > 0 {
				// Confirm selection from selection view
				paths := ms.getSelectedPaths()
				selectionType := ms.getSelectionType()
				return nil, FileSelection(paths, selectionType)
			}
			// Handle in filepicker view below

		case "delete":
			if ms.showingSelection {
				// Remove selected item from selection list
				if selected, ok := ms.selectedList.SelectedItem().(SelectedFileItem); ok {
					ms.removeSelection(selected.Path)
				}
				return ms, nil
			}

		case "space":
			if !ms.showingSelection {
				// Add current file to selection if in filepicker view
				if path := ms.filepicker.Path; path != "" {
					// Extract filename from path for display
					name := path
					if lastSlash := strings.LastIndex(path, "/"); lastSlash >= 0 {
						name = path[lastSlash+1:]
					}
					ms.addSelection(path, name)
				}
				return ms, nil
			}
		}
	}

	// Update the appropriate view
	if ms.showingSelection {
		var cmd tea.Cmd
		ms.selectedList, cmd = ms.selectedList.Update(msg)
		return ms, cmd
	} else {
		var cmd tea.Cmd
		ms.filepicker, cmd = ms.filepicker.Update(msg)

		// Check if user selected a file in filepicker
		if didSelect, path := ms.filepicker.DidSelectFile(msg); didSelect {
			// Extract filename from path for display
			name := path
			if lastSlash := strings.LastIndex(path, "/"); lastSlash >= 0 {
				name = path[lastSlash+1:]
			}
			ms.addSelection(path, name)
		}

		return ms, cmd
	}
}

func (ms MultiSelector) View() string {
	if ms.showingSelection {
		if len(ms.selected) == 0 {
			return ms.selectedList.View() + "\n\nNo files selected. Press Tab to return to file picker."
		}
		return ms.selectedList.View() + "\n\nPress Tab to return to file picker, Enter to confirm selection."
	}

	view := ms.filepicker.View()
	if len(ms.selected) > 0 {
		view += fmt.Sprintf("\n\nSelected: %d files (Press Tab to view, Space to add current file)", len(ms.selected))
	} else {
		view += "\n\nPress Space to select files, Tab to view selections"
	}
	return view
}
