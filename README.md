# Polymer

A composable Terminal User Interface (TUI) framework for Go, built on top of [Bubble Tea](https://github.com/charmbracelet/bubbletea). Polymer provides a simpler, more structured approach to building complex TUI applications with stack-based navigation, lifecycle hooks, and reusable components.

## Features

- **Simple Interface**: Minimal `Atom` interface with just `Name()`, `Update()`, and `View()` methods
- **Stack-based Navigation**: Push, pop, replace, and reset operations for intuitive UI flow
- **Lifecycle Hooks**: Debug and observe your application with `Lens` components
- **Composable Components**: Reusable UI elements like menus and forms
- **Bubble Tea Integration**: Seamless compatibility with existing Bubble Tea components
- **Type Safety**: Full Go type safety with generics support

## Installation

```bash
go get github.com/trippwill/polymer
```

## Quick Start

Here's a simple "Hello World" application:

```go
package main

import (
    "fmt"
    "os"
    
    tea "github.com/charmbracelet/bubbletea"
    poly "github.com/trippwill/polymer"
)

// HelloAtom is a simple screen that displays a greeting
type HelloAtom struct{}

func (h HelloAtom) Name() string { return "Hello" }

func (h HelloAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    // Quit on any key press
    if _, ok := msg.(tea.KeyMsg); ok {
        return h, tea.Quit
    }
    return h, nil
}

func (h HelloAtom) View() string {
    return "Hello, Polymer! (press any key to exit)\n"
}

func main() {
    // Create a host with your root atom
    host := poly.NewHost(HelloAtom{}, "Hello App")
    
    // Start the Bubble Tea program
    program := tea.NewProgram(host)
    if _, err := program.Run(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}
```

## Core Concepts

### Atom

An `Atom` is the fundamental building block of a Polymer application. It's similar to a Bubble Tea `Model` but with a simplified interface:

```go
type Atom interface {
    Name() string                           // Identifier for the atom
    Update(tea.Msg) (Atom, tea.Cmd)        // Handle messages and return new state
    View() string                           // Render the current state
}
```

Atoms can optionally implement the `Initializer` interface to provide setup logic:

```go
type Initializer interface {
    Atom
    Init() tea.Cmd
}
```

### Chain

A `Chain` manages a stack of Atoms, enabling complex navigation patterns:

```go
// Create a navigation chain starting with a menu
nav := poly.NewChain(mainMenu)

// Navigation commands
poly.Push(newAtom)     // Add atom to top of stack
poly.Pop()             // Remove top atom
poly.Replace(newAtom)  // Replace top atom
poly.Reset(newAtom)    // Clear stack and set new root
```

### Lens

A `Lens` provides lifecycle hooks for debugging and observability:

```go
// Create a lens with logging hooks
lens := poly.WithLens(myAtom,
    poly.WithOnInit(func(atom poly.Atom, cmd tea.Cmd) {
        log.Printf("Atom %s initialized", atom.Name())
    }),
    poly.WithBeforeUpdate(func(atom poly.Atom, msg tea.Msg) {
        log.Printf("Atom %s received message %T", atom.Name(), msg)
    }),
    poly.WithAfterUpdate(func(atom poly.Atom, cmd tea.Cmd) {
        log.Printf("Atom %s updated with command %T", atom.Name(), cmd)
    }),
    poly.WithOnView(func(atom poly.Atom, view string) {
        log.Printf("Atom %s rendered", atom.Name())
    }),
)
```

### Host

A `Host` bridges Polymer Atoms to Bubble Tea's `Model` interface:

```go
host := poly.NewHost(rootAtom, "My App")
program := tea.NewProgram(host)
```

The host automatically handles `Ctrl+C` and `q` for quitting the application.

## Examples

### Navigation Example

```go
// Create a menu with navigation
mainMenu := menu.NewMenu("Main Menu",
    menu.NewMenuItem(SettingsAtom{}, "Settings"),
    menu.NewMenuItem(AboutAtom{}, "About"),
    menu.NewMenuItem(QuitAtom{}, "Exit"),
)

// Wrap in a navigation chain
nav := poly.NewChain(mainMenu)

// The menu will automatically push selected atoms onto the stack
// Users can press 'esc' to pop back to previous screens
```

### Form Wizard Example

```go
// Multi-step form using stack navigation
type Step1 struct { /* ... */ }
type Step2 struct { /* ... */ }
type Step3 struct { /* ... */ }

func (s Step1) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            // Move to next step
            return s, poly.Push(Step2{data: s.collectData()})
        }
    }
    return s, nil
}
```

### Error Handling

```go
func (a MyAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case poly.ErrorMsg:
        // Handle errors from commands
        return ErrorAtom{error: error(msg)}, nil
    }
    return a, nil
}

// Emit errors from commands
func someFailingCommand() tea.Cmd {
    return poly.Error(fmt.Errorf("something went wrong"))
}
```

## Gels (Components)

Polymer includes reusable components called "gels":

### Menu

```go
import "github.com/trippwill/polymer/gels/menu"

// Create a menu with items
menu := menu.NewMenu("Choose an option",
    menu.NewMenuItem(atom1, "First Option"),
    menu.NewMenuItem(atom2, "Second Option"),
)
```

## Adapters

Polymer provides adapters for interoperability with Bubble Tea:

### Using Bubble Tea Models in Polymer

```go
// Wrap a Bubble Tea model as an Atom
bubbleTeaModel := textinput.New()
atom := poly.AtomAdapter{
    Model: bubbleTeaModel,
    AtomName: "Text Input",
}
```

### Using Polymer Atoms in Bubble Tea

```go
// Wrap a Polymer Atom as a Bubble Tea Model
atom := MyAtom{}
model := poly.WrapAtom(atom)
```

## Debugging and Tracing

Use the trace package for detailed logging:

```go
import "github.com/trippwill/polymer/trace"

// Set up logging
logger := log.New(os.Stderr, "debug: ", log.LstdFlags)

// Apply logging to your atom
traced := poly.WithLens(myAtom, trace.WithLogging(logger)...)
```

## Building and Running

```bash
# Build the library
go build ./...

# Run the example
cd examples/wizard
go run main.go

# The example demonstrates:
# - Menu navigation
# - Stack-based UI flow
# - Input handling
# - Multi-screen applications
```

## Architecture

Polymer follows these design principles:

1. **Composition over Inheritance**: Build complex UIs by composing simple Atoms
2. **Unidirectional Data Flow**: Messages flow down, state flows up
3. **Immutable State**: Atoms return new instances rather than mutating themselves
4. **Minimal Interface**: The `Atom` interface has just three methods
5. **Bubble Tea Compatibility**: Seamlessly integrate with the Bubble Tea ecosystem

## API Reference

### Core Types

- `Atom` - The main interface for UI components
- `Chain` - Stack-based navigation manager
- `Lens` - Lifecycle hooks and observability
- `Host` - Bubble Tea integration layer

### Navigation Messages

- `PushMsg` - Add atom to stack
- `PopMsg` - Remove top atom from stack  
- `ReplaceMsg` - Replace top atom
- `ResetMsg` - Reset stack to single atom

### Utility Functions

- `OptionalInit(atom)` - Call Init() if atom implements Initializer
- `Push(atom)`, `Pop()`, `Replace(atom)`, `Reset(atom)` - Navigation commands
- `Error(err)` - Convert error to command

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.