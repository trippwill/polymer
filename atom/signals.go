package atom

import tea "github.com/charmbracelet/bubbletea"

// ErrCode is failure in the [atom] package.
type ErrCode int

const (
	ErrUnkown       ErrCode = iota // Unknown error.
	ErrStackIsEmpty                // Stack is empty.
)

func (e ErrCode) Error() string {
	switch e {
	case ErrStackIsEmpty:
		return "atom: stack is empty"

	default:
		return "atom: unknown error"
	}
}

// Signals for [Stack] state changes.
type (
	PushMsg    struct{ State Model }
	PopMsg     struct{}
	ReplaceMsg struct{ State Model }
	ResetMsg   struct{ State Model }
)

// Push adds a new [Model] to the top of the [Stack], making it the [Stack.Active] state.
func Push(atom Model) tea.Cmd { return func() tea.Msg { return PushMsg{State: atom} } }

// Pop removes the top [Model] from the [Stack], returning it to the previous state.
func Pop() tea.Cmd { return func() tea.Msg { return PopMsg{} } }

// Replace replaces the current [Model] in the [Stack] with a new one, without changing the stack depth.
func Replace(atom Model) tea.Cmd { return func() tea.Msg { return ReplaceMsg{State: atom} } }

// Reset replaces the current [Model] in the [Stack] with a new one, resetting the stack to a single state.
// When atom is nil, it resets the stack to its initial state.
func Reset(atom Model) tea.Cmd { return func() tea.Msg { return ResetMsg{State: atom} } }
