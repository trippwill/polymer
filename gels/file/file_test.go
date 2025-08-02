package file

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

func TestNewSelector(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected func(*Selector) bool
	}{
		{
			name:   "default config",
			config: Config{},
			expected: func(s *Selector) bool {
				return s.name == "Select File" && s.config.CurrentDir != ""
			},
		},
		{
			name: "files only config",
			config: Config{
				Title:    "Pick File",
				FileType: FilesOnly,
			},
			expected: func(s *Selector) bool {
				return s.name == "Pick File" &&
					s.filepicker.FileAllowed && !s.filepicker.DirAllowed
			},
		},
		{
			name: "dirs only config",
			config: Config{
				Title:    "Pick Directory",
				FileType: DirsOnly,
			},
			expected: func(s *Selector) bool {
				return s.name == "Pick Directory" &&
					!s.filepicker.FileAllowed && s.filepicker.DirAllowed
			},
		},
		{
			name: "files and dirs config",
			config: Config{
				Title:      "Pick Anything",
				FileType:   FilesAndDirs,
				ShowHidden: true,
			},
			expected: func(s *Selector) bool {
				return s.name == "Pick Anything" &&
					s.filepicker.FileAllowed && s.filepicker.DirAllowed &&
					s.filepicker.ShowHidden
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewSelector(tt.config)
			if !tt.expected(selector) {
				t.Errorf("NewSelector() failed validation for %s", tt.name)
			}
		})
	}
}

func TestSelectorAtomInterface(t *testing.T) {
	selector := NewSelector(Config{Title: "Test"})

	// Test that it implements poly.Atom
	var _ atom.Model = *selector

	// Test Name method
	if name := selector.Name(); name != "Test" {
		t.Errorf("Name() = %q, want %q", name, "Test")
	}

	// Test Init method doesn't panic
	cmd := selector.Init()
	if cmd == nil {
		t.Error("Init() returned nil command")
	}

	// Test View method doesn't panic
	view := selector.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestSelectorUpdate(t *testing.T) {
	selector := NewSelector(Config{Title: "Test"})

	tests := []struct {
		name      string
		msg       tea.Msg
		expectPop bool
	}{
		{
			name:      "escape key",
			msg:       tea.KeyMsg{Type: tea.KeyEsc},
			expectPop: true,
		},
		{
			name:      "other key",
			msg:       tea.KeyMsg{Type: tea.KeyEnter},
			expectPop: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atom, cmd := selector.Update(tt.msg)

			// Should return same type
			if _, ok := atom.(Selector); !ok {
				t.Error("Update() didn't return Selector")
			}

			// For escape key, we can't easily test Pop() command
			// but we can verify cmd is not nil
			if tt.expectPop && cmd == nil {
				t.Error("Expected command for escape key")
			}
		})
	}
}

func TestNewMultiSelector(t *testing.T) {
	selector := NewMultiSelector(Config{
		Title:    "Multi Test",
		FileType: FilesAndDirs,
	})

	if selector.name != "Multi Test" {
		t.Errorf("Name = %q, want %q", selector.name, "Multi Test")
	}

	if selector.showingSelection {
		t.Error("Should start in filepicker view")
	}

	if len(selector.selected) != 0 {
		t.Error("Should start with no selections")
	}
}

func TestMultiSelectorAtomInterface(t *testing.T) {
	selector := NewMultiSelector(Config{Title: "Multi Test"})

	// Test that it implements poly.Atom
	var _ atom.Model = *selector

	// Test Name method
	if name := selector.Name(); name != "Multi Test" {
		t.Errorf("Name() = %q, want %q", name, "Multi Test")
	}

	// Test Init method doesn't panic
	cmd := selector.Init()
	if cmd == nil {
		t.Error("Init() returned nil command")
	}

	// Test View method doesn't panic
	view := selector.View()
	if view == "" {
		t.Error("View() returned empty string")
	}
}

func TestMultiSelectorHelperMethods(t *testing.T) {
	selector := NewMultiSelector(Config{Title: "Test"})

	// Test addSelection
	selector.addSelection("/test/file.txt", "file.txt")

	if len(selector.selected) != 1 {
		t.Errorf("After addSelection, len(selected) = %d, want 1", len(selector.selected))
	}

	if item, exists := selector.selected["/test/file.txt"]; !exists {
		t.Error("File not found in selected map")
	} else if item.Name != "file.txt" || item.Path != "/test/file.txt" {
		t.Errorf("Selected item = %+v, want Name=file.txt Path=/test/file.txt", item)
	}

	// Test getSelectedPaths
	paths := selector.getSelectedPaths()
	if len(paths) != 1 || paths[0] != "/test/file.txt" {
		t.Errorf("getSelectedPaths() = %v, want [/test/file.txt]", paths)
	}

	// Test removeSelection
	selector.removeSelection("/test/file.txt")

	if len(selector.selected) != 0 {
		t.Errorf("After removeSelection, len(selected) = %d, want 0", len(selector.selected))
	}
}

func TestMultiSelectorUpdate(t *testing.T) {
	selector := NewMultiSelector(Config{Title: "Test"})

	tests := []struct {
		name                     string
		msg                      tea.Msg
		expectedShowingSelection bool
	}{
		{
			name:                     "tab key toggles view",
			msg:                      tea.KeyMsg{Type: tea.KeyTab},
			expectedShowingSelection: true,
		},
		{
			name:                     "escape in selection view goes back to filepicker",
			msg:                      tea.KeyMsg{Type: tea.KeyEsc},
			expectedShowingSelection: false,
		},
	}

	currentSelector := *selector // Work with value type for testing
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atom, cmd := currentSelector.Update(tt.msg)

			updatedSelector, ok := atom.(MultiSelector)
			if !ok {
				t.Fatal("Update() didn't return MultiSelector")
			}

			if updatedSelector.showingSelection != tt.expectedShowingSelection {
				t.Errorf("showingSelection = %v, want %v",
					updatedSelector.showingSelection, tt.expectedShowingSelection)
			}

			// Tab should not generate a command
			if tt.name == "tab key toggles view" && cmd != nil {
				t.Error("Tab key should not generate command")
			}

			// Update selector for next test
			currentSelector = updatedSelector
		})
	}
}

func TestSelectionType(t *testing.T) {
	// Create temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "polymer_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	testFile := filepath.Join(tmpDir, "test.txt")
	testDir := filepath.Join(tmpDir, "testdir")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		config     Config
		selections map[string]SelectedFileItem
		expected   SelectionType
	}{
		{
			name:   "single file",
			config: Config{FileType: FilesOnly},
			selections: map[string]SelectedFileItem{
				testFile: {Name: "test.txt", Path: testFile},
			},
			expected: SelectionTypeFile,
		},
		{
			name:   "multiple files",
			config: Config{FileType: FilesOnly},
			selections: map[string]SelectedFileItem{
				testFile:       {Name: "test.txt", Path: testFile},
				testFile + "2": {Name: "test2.txt", Path: testFile + "2"},
			},
			expected: SelectionTypeFiles,
		},
		{
			name:   "single directory",
			config: Config{FileType: DirsOnly},
			selections: map[string]SelectedFileItem{
				testDir: {Name: "testdir", Path: testDir},
			},
			expected: SelectionTypeDirectory,
		},
		{
			name:   "multiple directories",
			config: Config{FileType: DirsOnly},
			selections: map[string]SelectedFileItem{
				testDir:       {Name: "testdir", Path: testDir},
				testDir + "2": {Name: "testdir2", Path: testDir + "2"},
			},
			expected: SelectionTypeDirectories,
		},
		{
			name:   "mixed selection",
			config: Config{FileType: FilesAndDirs},
			selections: map[string]SelectedFileItem{
				testFile: {Name: "test.txt", Path: testFile},
				testDir:  {Name: "testdir", Path: testDir},
			},
			expected: SelectionTypeMixed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewMultiSelector(tt.config)
			selector.selected = tt.selections

			selectionType := selector.getSelectionType()
			if selectionType != tt.expected {
				t.Errorf("getSelectionType() = %v, want %v", selectionType, tt.expected)
			}
		})
	}
}

func TestFileSelectionMsg(t *testing.T) {
	// Test the FileSelectionMsg and FileSelection command
	files := []string{"/test/file1.txt", "/test/file2.txt"}
	selectionType := SelectionTypeFiles

	cmd := FileSelection(files, selectionType)
	if cmd == nil {
		t.Fatal("FileSelection returned nil command")
	}

	// Execute the command to get the message
	msg := cmd()

	fileMsg, ok := msg.(FileSelectionMsg)
	if !ok {
		t.Fatalf("Command returned %T, want FileSelectionMsg", msg)
	}

	if len(fileMsg.Files) != 2 {
		t.Errorf("Files length = %d, want 2", len(fileMsg.Files))
	}

	if fileMsg.Type != selectionType {
		t.Errorf("Type = %v, want %v", fileMsg.Type, selectionType)
	}

	for i, file := range files {
		if fileMsg.Files[i] != file {
			t.Errorf("Files[%d] = %q, want %q", i, fileMsg.Files[i], file)
		}
	}
}

