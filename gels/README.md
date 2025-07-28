# Polymer Gels

Gels are reusable UI components for Polymer applications. They implement the `Atom` interface and can be easily composed into larger applications.

## Menu Gel

**Package**: `github.com/trippwill/polymer/gels/menu`

A customizable menu component built on Bubble Tea's `list` component, with automatic navigation integration.

### Features

- Keyboard navigation (up/down arrows, j/k)
- Item selection with Enter
- Back navigation with Esc
- Search/filtering with `/`
- Automatic window resizing
- Customizable items and descriptions

### Usage

```go
import "github.com/trippwill/polymer/gels/menu"

// Create menu items
items := []menu.Item{
    menu.NewMenuItem(SettingsAtom{}, "Application Settings"),
    menu.NewMenuItem(HelpAtom{}, "Help and Documentation"), 
    menu.NewMenuItem(QuitAtom{}, "Exit Application"),
}

// Create the menu
mainMenu := menu.NewMenu("Main Menu", items...)

// Use in a Chain for automatic navigation
nav := poly.NewChain(mainMenu)
```

### MenuItem Structure

```go
type Item struct {
    Atom        poly.Atom  // The atom to navigate to when selected
    Description string     // Descriptive text shown in the menu
}

// Menu items implement list.DefaultItem for Bubble Tea compatibility
func (m Item) Title() string       { return m.Atom.Name() }
func (m Item) Description() string { return m.description }
func (m Item) FilterValue() string { return title + " " + description }
```

### Key Bindings

- `↑/k` - Move selection up
- `↓/j` - Move selection down  
- `/` - Start filtering/search
- `Enter` - Select item (pushes atom onto navigation stack)
- `Esc` - Go back (pops from navigation stack)
- `q` - Quit application (handled by Host)

### Customization

The menu uses Bubble Tea's default list delegate, but you can customize appearance by modifying the list after creation:

```go
menu := menu.NewMenu("Custom Menu", items...)

// Access the underlying list for customization
menu.List().SetShowStatusBar(false)
menu.List().SetShowPagination(false)
menu.List().SetShowHelp(false)
```

### Integration Example

```go
// Menu automatically integrates with Chain navigation
type SettingsAtom struct {
    // Settings implementation
}

func (s SettingsAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "esc" {
            return s, poly.Pop() // Return to previous menu
        }
    }
    return s, nil
}

// When user selects "Settings" from menu, SettingsAtom is automatically pushed
```

## Creating Custom Gels

To create your own reusable gel:

1. **Implement the Atom interface**:
```go
type MyGel struct {
    // State fields
}

func (g MyGel) Name() string { return "MyGel" }
func (g MyGel) Update(msg tea.Msg) (poly.Atom, tea.Cmd) { /* implementation */ }
func (g MyGel) View() string { /* implementation */ }
```

2. **Provide a constructor function**:
```go
func NewMyGel(config MyGelConfig) *MyGel {
    return &MyGel{
        // Initialize with config
    }
}
```

3. **Add navigation support if needed**:
```go
func (g MyGel) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "enter":
            return g, poly.Push(g.selectedItem)
        case "esc": 
            return g, poly.Pop()
        }
    }
    return g, nil
}
```

4. **Make it composable**:
```go
// Allow embedding other atoms
type MyGel struct {
    items []poly.Atom
    // other fields
}
```

## Planned Gels

Future gel components may include:

- **Form**: Input forms with validation
- **Table**: Data tables with sorting and filtering
- **Tabs**: Tabbed interface component
- **Dialog**: Modal dialogs and confirmations
- **Progress**: Progress bars and spinners
- **Tree**: Hierarchical tree navigation
- **Split**: Split pane layouts

## Contributing Gels

When contributing new gels:

1. Place them in `gels/<name>/` directory
2. Follow the existing naming patterns
3. Include comprehensive documentation
4. Add example usage in the `examples/` directory
5. Ensure they work well with Chain navigation
6. Add tests for complex logic