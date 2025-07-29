# File Selector Gels

File selector gels provide file and directory selection capabilities for Polymer applications. They support both single and multiple selection modes with full keyboard navigation.

## Features

- **File System Navigation**: Browse directories, navigate up/down the directory tree
- **Configurable Selection**: Select files only, directories only, or both
- **Single Selection**: Select one file or directory at a time
- **Multi-Selection**: Select multiple files/directories with visual feedback
- **Hidden Files**: Configurable display of hidden files (starting with '.')
- **Keyboard Navigation**: Full keyboard support with intuitive key bindings
- **Visual Feedback**: Clear indication of selected items and current directory
- **Remove from Selection**: Ability to remove items from multi-selection

## Single File Selector

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

// Use in a Chain for navigation
nav := poly.NewChain(selector)
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

- `↑/↓` or `j/k` - Navigate up/down in file list
- `/` - Start filtering/search
- `Enter` - Select file or enter directory
- `Esc` - Go back to previous screen
- `Backspace` - Navigate to parent directory
- `q` - Quit application

## Multi-File Selector

### Usage

```go
import "github.com/trippwill/polymer/gels/file"

// Create a multi-file selector
multiSelector := file.NewMultiSelector(file.Config{
    Title:      "Select Multiple Files",
    FileType:   file.FilesOnly,
    ShowHidden: false,
})
```

### Key Bindings

All single selector keys plus:

- `Space` - Toggle selection of current item
- `Ctrl+A` - Select all items in current directory
- `Ctrl+D` - Deselect all items
- `Delete` - Remove currently highlighted item from selection
- `Enter` - Complete selection (when items are selected) or navigate into directory

### Selection Feedback

The multi-selector provides clear visual feedback:

- **Title Updates**: Shows count of selected items: `"Select Files (3 selected)"`
- **Selection Summary**: Displays selected filenames at bottom of screen
- **Current Directory**: Always shows current path in title

## Example

```go
package main

import (
    "fmt"
    "log"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    poly "github.com/trippwill/polymer"
    "github.com/trippwill/polymer/gels/file"
    "github.com/trippwill/polymer/gels/menu"
)

func main() {
    // Single file selector
    singleFile := file.NewSelector(file.Config{
        Title:    "Select Single File",
        FileType: file.FilesOnly,
    })

    // Multi-file selector
    multiFile := file.NewMultiSelector(file.Config{
        Title:    "Select Multiple Files",
        FileType: file.FilesAndDirs,
        ShowHidden: true,
    })

    // Directory selector
    dirSelector := file.NewSelector(file.Config{
        Title:    "Select Directory",
        FileType: file.DirsOnly,
    })

    // Create menu
    mainMenu := menu.NewMenu(
        "File Selector Demo",
        menu.NewMenuItem(singleFile, "Single File Selection"),
        menu.NewMenuItem(multiFile, "Multiple File Selection"),
        menu.NewMenuItem(dirSelector, "Directory Selection"),
    )

    // Run application
    host := poly.NewHost("File Selector", poly.NewChain(mainMenu))
    program := tea.NewProgram(host)
    program.Run()
}
```

## Integration with Polymer

File selectors integrate seamlessly with Polymer's navigation system:

- **Chain Navigation**: Use `poly.Push()` to add selectors to the navigation stack
- **Notifications**: Selection results are sent via `poly.Notify()`
- **Lifecycle Hooks**: Support all Polymer lens features for debugging
- **Error Handling**: Graceful handling of file system errors

## File Selection Results

When a selection is made, the selector:

1. Sends a notification with selected file paths
2. Automatically pops back to the previous screen
3. Returns the full path(s) of selected items

For multi-selection, the notification includes all selected file paths in a comma-separated list.

## Customization

The file selector uses Bubble Tea's list component and can be customized:

```go
selector := file.NewSelector(config)

// Access underlying list for customization
// Note: This would require exposing the list field or adding methods
```

Both single and multi-selectors follow Polymer's architectural patterns and can be easily extended or customized for specific use cases.