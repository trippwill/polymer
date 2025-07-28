# Polymer Examples

This directory contains example applications demonstrating various Polymer features and patterns.

## Wizard Example

**Location**: `./wizard/`

A comprehensive example showcasing:
- Menu-based navigation using the `menu` gel
- Stack-based UI flow with push/pop operations
- Input handling and form-like interactions
- Multi-screen application architecture
- Debugging with trace logging

### Running the Wizard

```bash
cd wizard
go run main.go
```

### What it demonstrates

1. **Main Menu**: Uses the `menu.NewMenu()` gel to create a selectable list
2. **Name Input Screen**: Demonstrates text input handling and validation
3. **Greeting Screen**: Shows dynamic content based on user input
4. **Navigation Flow**: 
   - Menu → Name Input (push)
   - Name Input → Greeting (push) 
   - Greeting → Menu (double pop)
5. **Lifecycle Hooks**: Logs all Atom lifecycle events to `debug.log`

### Key Code Patterns

```go
// Menu creation with navigation
mainMenu := menu.NewMenu("Main Menu",
    menu.NewMenuItem(NamePromptScreen{}, "Run the Name Wizard"),
    menu.NewMenuItem(QuitAtom{}, "Exit Application"),
)

// Stack navigation
poly.Push(GreetingScreen{Value: n.input})  // Move forward
tea.Batch(poly.Pop(), poly.Pop())          // Go back multiple steps

// Tracing and debugging  
traced := poly.WithLens(nav, trace.WithLogging(logger)...)
```

### Architecture

```
Host (Bubble Tea Model)
└── Lens (with logging hooks)
    └── Chain (navigation stack)
        └── Menu → NamePrompt → Greeting
```

The example shows how to build a complete application with multiple screens, user input, and navigation - all with minimal boilerplate code.

## Adding New Examples

To create a new example:

1. Create a new directory under `examples/`
2. Add a `main.go` file with your Polymer application
3. Update this README with a description of your example
4. Follow the existing patterns for consistency

Good example ideas:
- Form validation
- Data tables/lists
- File browser
- Settings screens
- Multi-pane layouts
- Real-time updates