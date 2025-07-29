package polymer

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// Chain manages a stack of Atoms.
type Chain struct {
	stack []Atom
	name  string
}

// Messages for the Chain to handle state changes.
type (
	PushMsg    struct{ State Atom }
	PopMsg     struct{}
	ReplaceMsg struct{ State Atom }
	ResetMsg   struct{ State Atom }
)

var ErrChainEmpty = fmt.Errorf("chain is empty")

// PushMsg, PopMsg, ReplaceMsg, and ResetMsg are used to manipulate the Chain's stack.
func Push(atom Atom) tea.Cmd    { return func() tea.Msg { return PushMsg{State: atom} } }
func Pop() tea.Cmd              { return func() tea.Msg { return PopMsg{} } }
func Replace(atom Atom) tea.Cmd { return func() tea.Msg { return ReplaceMsg{State: atom} } }
func Reset(atom Atom) tea.Cmd   { return func() tea.Msg { return ResetMsg{State: atom} } }

// NewChain creates a new Chain with an initial Atom.
// It panics if the initial Atom is nil.
func NewChain(initial Atom) *Chain {
	if initial == nil {
		panic("initial state cannot be nil")
	}

	return &Chain{
		stack: []Atom{initial},
	}
}

// Active returns the Atom at the top of the stack.
func (h *Chain) Active() Atom {
	if len(h.stack) == 0 {
		return nil
	}
	return h.stack[len(h.stack)-1]
}

// Peek returns the Atom at the specified position from the top of the stack.
// If n is out of bounds, it returns nil.
func (h *Chain) Peek(n int) Atom {
	if n < 0 || n >= len(h.stack) {
		return nil
	}
	return h.stack[len(h.stack)-1-n]
}

// Depth returns the number of Atoms in the stack.
func (h *Chain) Depth() int {
	return len(h.stack)
}

func (h *Chain) Name() string { return h.name }

// Init initializes the active state in the stack if it implements Initializer.
func (h *Chain) Init() tea.Cmd {
	current := h.Active()
	if initializer, ok := current.(interface{ Init() tea.Cmd }); ok {
		return initializer.Init()
	}
	return nil
}

// Update processes the message and updates the active state in the stack.
func (h *Chain) Update(msg tea.Msg) (Atom, tea.Cmd) {
	switch m := msg.(type) {
	case PushMsg:
		h.push(m.State)
		return h, OptionalInit(h.Active())
	case ReplaceMsg:
		h.replace(m.State)
		return h, OptionalInit(h.Active())
	case ResetMsg:
		h.reset(m.State)
		return h, OptionalInit(h.Active())
	case PopMsg:
		h.pop()
		return h, OptionalInit(h.Active())
	}

	current := h.Active()
	if current == nil {
		return h, Error(ErrChainEmpty)
	}

	next, cmd := current.Update(msg)
	h.replace(next)
	return h, cmd
}

// View returns the view of the active state in the stack.
func (h *Chain) View() string {
	current := h.Active()
	if current == nil {
		return ""
	}
	return current.View()
}

// push adds a new Atom to the top of the stack.
func (h *Chain) push(state Atom) {
	if state == nil {
		return
	}

	h.stack = append(h.stack, state)
}

// pop removes the top Atom from the stack and returns it.
func (h *Chain) pop() Atom {
	if len(h.stack) == 0 {
		return nil
	}

	active := h.Active()
	h.stack = h.stack[:len(h.stack)-1]
	return active
}

// replace the top Atom in the stack with a new one.
func (h *Chain) replace(atom Atom) {
	if atom == nil {
		return
	}

	if len(h.stack) > 0 {
		h.stack[len(h.stack)-1] = atom
	} else {
		h.stack = append(h.stack, atom)
	}
}

// reset the stack to a single Atom state.
// if atom is nil, it resets to the initial Atom in the stack.
func (h *Chain) reset(atom Atom) {
	if atom == nil {
		h.stack = []Atom{h.stack[0]}
		return
	}

	h.stack = []Atom{atom}
}
