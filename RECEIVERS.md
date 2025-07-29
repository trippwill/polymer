# Receiver Patterns in Polymer

This document explains the subtle but important distinctions around method receivers in Polymer, which builds on the BubbleTea framework.

## Core Principles

### Interface Methods (Model/Atom)

These methods **MUST** use value receivers and **NEVER** mutate state:

```go
// ✅ Correct: Value receiver, no mutation
func (m Model) Name() string { return m.name }
func (m Model) Init() tea.Cmd { return m.someCmd }
func (m Model) View() string { return m.render() }
```

```go
// ❌ Wrong: Mutating state in interface methods
func (m Model) Name() string { 
    m.counter++ // Never do this!
    return m.name 
}
```

### Update Method

The `Update` method should use a **value receiver** but **MAY** mutate the receiver within the method, as long as the updated value is returned:

```go
// ✅ Correct: Value receiver, mutation allowed, return updated value
func (m Model) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        m.counter++ // This is OK in Update()
        return m, nil
    }
    return m, nil
}
```

### Helper Methods (Unexported)

Unexported helper methods **MAY** use pointer receivers for efficiency, especially when:
- They modify large structs
- They're called frequently
- They need to mutate multiple fields

```go
// ✅ Efficient: Pointer receiver for helper methods
func (m *Model) updateInternalState() {
    m.items = append(m.items, newItem)
    m.list.SetItems(m.items)
    m.count = len(m.items)
}

// ✅ Also valid: Value receiver if returning modified copy
func (m Model) addItem(item Item) Model {
    m.items = append(m.items, item)
    return m
}
```

## Why This Matters

1. **Interface Contract**: BubbleTea expects interface methods to be pure/non-mutating
2. **Performance**: Pointer receivers avoid copying large structs for helper methods
3. **Predictability**: Clear separation between mutation points and pure functions
4. **Consistency**: Following established patterns makes code more maintainable

## Example: File Selector Implementation

```go
// Interface methods: value receivers, no mutation
func (s Selector) Name() string { return s.name }
func (s Selector) View() string { return s.filepicker.View() }

// Update: value receiver, mutation allowed within method
func (s Selector) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
    s.filepicker, cmd = s.filepicker.Update(msg) // Mutation OK
    return s, cmd // Return updated value
}

// Helper methods: pointer receivers for efficiency
func (ms *MultiSelector) addSelection(path, name string) {
    ms.selected[path] = SelectedFileItem{Name: name, Path: path}
    ms.updateSelectedList()
}
```

## Migration Guide

When refactoring from all-value receivers to this pattern:

1. Keep interface methods as value receivers
2. Keep Update as value receiver
3. Change helper methods to pointer receivers if they mutate state
4. Update calls to helper methods (remove assignment if using pointer receivers)

```go
// Before: value receiver helper
ms = ms.addSelection(path, name)

// After: pointer receiver helper  
ms.addSelection(path, name)
```

This pattern provides the best balance of safety, performance, and BubbleTea compatibility.