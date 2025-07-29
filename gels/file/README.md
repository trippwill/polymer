# File Selector Gels

File selector gels provide file and directory selection capabilities for Polymer applications, built on top of the Bubble Tea ecosystem's filepicker component. They support both single and multiple selection modes with full keyboard navigation.

## Features

- **Built on Bubbles**: Uses the well-tested `bubbles/filepicker` component
- **File System Navigation**: Browse directories, navigate up/down the directory tree
- **Configurable Selection**: Select files only, directories only, or both
- **Single Selection**: Select one file or directory at a time using bubbles/filepicker
- **Multi-Selection**: Select multiple files/directories with dual-view interface
- **Hidden Files**: Configurable display of hidden files (starting with '.')
- **Keyboard Navigation**: Full keyboard support with intuitive key bindings
- **Custom Messages**: Uses proper message passing instead of notifications
- **Visual Feedback**: Clear indication of selected items and current directory

## Screenshots

### Main Menu
![Main Menu](https://github.com/user-attachments/assets/5736cf03-4c67-43ed-9c76-d7bba31b7bbb)

The main menu shows all available file selector options, allowing users to choose between single file selection, directory selection, multi-file selection, and mixed selection modes.

## Single File Selector

The single file selector wraps `bubbles/filepicker` for Polymer integration.

### Screenshot
![Single File Selector](https://github.com/user-attachments/assets/3bb05f51-5d23-4939-bcc7-80e9b446882a)

The single file selector provides a clean interface for browsing directories and selecting individual files. It shows directories with folder icons (📁) and files with document icons (📄), including support for hidden files when enabled.

### Usage

```go
import "github.com/trippwill/polymer/gels/file"

// Create a single file selector
selector := file.NewSelector(file.Config{
    Title:      "Select File",
    FileType:   file.FilesOnly,    // FilesOnly, DirsOnly, or FilesAndDirs
    CurrentDir: "/path/to/start",  // Optional, defaults to current working directory
    ShowHidden: false,             // Whether to show hidden files
})

// Handle selection results
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case poly.FileSelectionMsg:
        // Handle file selection
        fmt.Printf("Selected %s: %s\n", msg.Type, msg.Files[0])
    }
    return m, nil
}
```

### Configuration Options

```go
type Config struct {
    Title       string    // Title displayed in the selector
    FileType    FileType  // What can be selected
    CurrentDir  string    // Starting directory (default: current working dir)
    ShowHidden  bool      // Show hidden files (default: false)
}

type FileType int

const (
    FilesOnly    FileType = iota  // Only files can be selected
    DirsOnly                      // Only directories can be selected  
    FilesAndDirs                  // Both files and directories
)
```

### Key Bindings

Standard bubbles/filepicker key bindings:

- `↑/↓` or `j/k` - Navigate up/down in file list
- `Enter` - Select file or enter directory  
- `Esc` - Go back to previous screen
- `/` - Start filtering/search (if enabled)
- `h/l` - Navigate to parent/child directory
- `q` - Quit application

## Multi-File Selector

The multi-file selector combines `bubbles/filepicker` with `bubbles/list` to provide a dual-view interface for selecting multiple files.

### Screenshots

#### File Picker View
![Multi-File Picker View](https://github.com/user-attachments/assets/eaaea596-9fba-40ef-9647-7c53f45f27f3)

The file picker view allows users to browse directories and add files to their selection using the Space bar. Selected files are marked with checkmarks (✓) and the current selection count is displayed.

#### Selection List View  
![Selection List View](https://github.com/user-attachments/assets/8bffbbdf-7951-49ce-910a-76f1e1c573e8)

The selection list view shows all currently selected files and allows users to review and manage their selections. Users can remove items with the Delete key or confirm their selection with Enter.

### Usage

```go
import "github.com/trippwill/polymer/gels/file"

// Create a multi-file selector
multiSelector := file.NewMultiSelector(file.Config{
    Title:      "Select Multiple Files",
    FileType:   file.FilesOnly,
    ShowHidden: false,
})

// Handle selection results
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    switch msg := msg.(type) {
    case poly.FileSelectionMsg:
        // Handle multiple file selection
        fmt.Printf("Selected %d %ss: %v\n", len(msg.Files), msg.Type, msg.Files)
    }
    return m, nil
}
```

### Key Bindings

Filepicker view:
- All standard filepicker keys plus:
- `Space` - Add current file to selection
- `Tab` - Switch to selection list view
- `Esc` - Complete selection (if files selected) or go back

Selection list view:
- `↑/↓` - Navigate selection list
- `Delete` - Remove item from selection
- `Tab` - Switch back to filepicker view
- `Enter` - Confirm selection
- `Esc` - Switch back to filepicker view

### Selection Interface

The multi-selector provides two views:

1. **Filepicker View**: Browse and add files using the standard filepicker
2. **Selection View**: Review and manage selected files using a list

Switch between views with `Tab`. The bottom of the screen shows selection status and available actions.

## Example

```go
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
)

// Handler that displays file selection results
type FileHandler struct {
    *menu.Model
    lastSelection string
}

func (h *FileHandler) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case poly.FileSelectionMsg:
        if len(msg.Files) == 1 {
            h.lastSelection = fmt.Sprintf("Selected %s: %s", msg.Type, msg.Files[0])
        } else {
            h.lastSelection = fmt.Sprintf("Selected %d files: %s", 
                len(msg.Files), strings.Join(msg.Files, ", "))
        }
        return h, nil
    }
    return h.Model.Update(msg)
}

func (h *FileHandler) View() string {
    view := h.Model.View()
    if h.lastSelection != "" {
        view += "\n\n" + h.lastSelection
    }
    return view
}

func main() {
    // Single file selector
    singleFile := file.NewSelector(file.Config{
        Title:    "Select Single File",
        FileType: file.FilesOnly,
    })

    // Multi-file selector
    multiFile := file.NewMultiSelector(file.Config{
        Title:      "Select Multiple Files",
        FileType:   file.FilesAndDirs,
        ShowHidden: true,
    })

    // Directory selector
    dirSelector := file.NewSelector(file.Config{
        Title:    "Select Directory",
        FileType: file.DirsOnly,
    })

    // Create menu with selection handler
    mainMenu := &FileHandler{
        Model: menu.NewMenu(
            "File Selector Demo",
            menu.NewMenuItem(singleFile, "Single File Selection"),
            menu.NewMenuItem(multiFile, "Multiple File Selection"),
            menu.NewMenuItem(dirSelector, "Directory Selection"),
        ),
    }

    // Run application
    host := poly.NewHost("File Selector", poly.NewChain(mainMenu))
    program := tea.NewProgram(host)
    program.Run()
}
```

## Integration with Polymer

File selectors integrate seamlessly with Polymer's navigation system:

- **Chain Navigation**: Use `poly.Push()` to add selectors to the navigation stack
- **Custom Messages**: Selection results are sent via `poly.FileSelectionMsg`
- **Lifecycle Hooks**: Support all Polymer lens features for debugging
- **Error Handling**: Graceful handling of file system errors via bubbles/filepicker

## Message Types

File selectors use custom message types for proper data flow:

```go
// FileSelectionMsg represents a file selection result
type FileSelectionMsg struct {
    Files []string  // Array of selected file paths
    Type  string    // "file", "directory", or "files"
}

// Create selection message
poly.FileSelection([]string{"/path/to/file"}, "file")
```

## Architecture

- **Single Selector**: Wraps `bubbles/filepicker` directly
- **Multi Selector**: Combines `bubbles/filepicker` + `bubbles/list`
- **Polymer Integration**: Custom Atom implementations with proper message passing
- **Configuration**: Shared Config struct for consistent API

This approach follows Polymer's principle of building on proven Bubble Tea ecosystem components while providing clean integration and enhanced functionality.