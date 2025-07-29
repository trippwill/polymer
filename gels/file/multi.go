package file

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// MultiSelector extends Selector to support multiple selections
type MultiSelector struct {
	*Selector
	selected    map[string]FileItem // map of path -> FileItem for selected items
	showingHelp bool
}

// NewMultiSelector creates a new multi-file selector
func NewMultiSelector(config Config) *MultiSelector {
	if config.Title == "" {
		config.Title = "Select Files"
	}

	selector := NewSelector(config)
	ms := &MultiSelector{
		Selector: selector,
		selected: make(map[string]FileItem),
	}

	ms.name = config.Title
	ms.setupMultiList()
	return ms
}

func (ms *MultiSelector) setupMultiList() {
	ms.setupList() // Call parent setup
	
	// Override help to include multi-selection keys
	ms.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "navigate/confirm"),
			),
			key.NewBinding(
				key.WithKeys("space"),
				key.WithHelp("space", "toggle selection"),
			),
			key.NewBinding(
				key.WithKeys("ctrl+a"),
				key.WithHelp("ctrl+a", "select all"),
			),
			key.NewBinding(
				key.WithKeys("ctrl+d"),
				key.WithHelp("ctrl+d", "deselect all"),
			),
			key.NewBinding(
				key.WithKeys("delete"),
				key.WithHelp("del", "remove selected"),
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

	// Update title to show selection count
	ms.updateTitle()
}

func (ms *MultiSelector) updateTitle() {
	baseTitle := fmt.Sprintf("%s - %s", ms.config.Title, ms.currentDir)
	if len(ms.selected) > 0 {
		ms.list.Title = fmt.Sprintf("%s (%d selected)", baseTitle, len(ms.selected))
	} else {
		ms.list.Title = baseTitle
	}
}

func (ms *MultiSelector) isSelected(item FileItem) bool {
	_, exists := ms.selected[item.Path]
	return exists
}

func (ms *MultiSelector) toggleSelection(item FileItem) {
	if ms.isSelected(item) {
		delete(ms.selected, item.Path)
	} else {
		ms.selected[item.Path] = item
	}
	ms.updateTitle()
}

func (ms *MultiSelector) selectAll() {
	for _, item := range ms.list.Items() {
		if fileItem, ok := item.(FileItem); ok {
			// Only select files/dirs that match our filter criteria and aren't parent dir
			if fileItem.Name != ".." {
				switch ms.config.FileType {
				case FilesOnly:
					if !fileItem.IsDir {
						ms.selected[fileItem.Path] = fileItem
					}
				case DirsOnly:
					if fileItem.IsDir {
						ms.selected[fileItem.Path] = fileItem
					}
				case FilesAndDirs:
					ms.selected[fileItem.Path] = fileItem
				}
			}
		}
	}
	ms.updateTitle()
}

func (ms *MultiSelector) deselectAll() {
	ms.selected = make(map[string]FileItem)
	ms.updateTitle()
}

func (ms *MultiSelector) getSelectedList() []string {
	var selected []string
	for path := range ms.selected {
		selected = append(selected, path)
	}
	return selected
}

var _ poly.Atom = (*MultiSelector)(nil)

func (ms *MultiSelector) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ms.width = msg.Width
		ms.height = msg.Height
		ms.list.SetSize(msg.Width, msg.Height)
		return ms, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "space":
			if selected, ok := ms.list.SelectedItem().(FileItem); ok {
				// Don't allow selecting parent directory
				if selected.Name != ".." {
					ms.toggleSelection(selected)
				}
			}
			return ms, nil

		case "ctrl+a":
			ms.selectAll()
			return ms, nil

		case "ctrl+d":
			ms.deselectAll()
			return ms, nil

		case "delete":
			if selected, ok := ms.list.SelectedItem().(FileItem); ok {
				if ms.isSelected(selected) {
					delete(ms.selected, selected.Path)
					ms.updateTitle()
				}
			}
			return ms, nil

		case "enter":
			if selected, ok := ms.list.SelectedItem().(FileItem); ok {
				if selected.IsDir {
					// Navigate into directory
					ms.currentDir = selected.Path
					ms.Selector.currentDir = selected.Path // Keep parent in sync
					ms.setupMultiList()
					if ms.width > 0 && ms.height > 0 {
						ms.list.SetSize(ms.width, ms.height)
					}
					return ms, nil
				} else if len(ms.selected) > 0 {
					// Complete selection with current selections
					selectedPaths := ms.getSelectedList()
					return ms, tea.Sequence(
						poly.Notify(poly.InfoLevel, fmt.Sprintf("Selected %d files: %s", len(selectedPaths), strings.Join(selectedPaths, ", "))),
						poly.Pop(),
					)
				} else {
					// Add current file to selection if none selected
					ms.toggleSelection(selected)
					return ms, nil
				}
			}

		case "backspace":
			// Navigate to parent directory
			parent := filepath.Dir(ms.currentDir)
			if parent != ms.currentDir { // Avoid infinite loop at root
				ms.currentDir = parent
				ms.Selector.currentDir = parent // Keep parent in sync
				ms.setupMultiList()
				if ms.width > 0 && ms.height > 0 {
					ms.list.SetSize(ms.width, ms.height)
				}
				return ms, nil
			}

		case "esc":
			if len(ms.selected) > 0 {
				// Show confirmation or complete selection
				selectedPaths := ms.getSelectedList()
				return ms, tea.Sequence(
					poly.Notify(poly.InfoLevel, fmt.Sprintf("Selected %d files: %s", len(selectedPaths), strings.Join(selectedPaths, ", "))),
					poly.Pop(),
				)
			}
			return ms, poly.Pop()
		}
	}

	var cmd tea.Cmd
	ms.list, cmd = ms.list.Update(msg)
	return ms, cmd
}

func (ms *MultiSelector) View() string {
	baseView := ms.list.View()
	
	// Add selection summary at the bottom if there are selections
	if len(ms.selected) > 0 {
		summary := fmt.Sprintf("\n[Selected: %d items]", len(ms.selected))
		if len(ms.selected) <= 3 {
			var names []string
			for _, item := range ms.selected {
				names = append(names, item.Name)
			}
			summary += fmt.Sprintf(" %s", strings.Join(names, ", "))
		}
		baseView += summary
	}
	
	return baseView
}