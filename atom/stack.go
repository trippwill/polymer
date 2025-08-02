package atom

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// Stack manages a stack of atom.Models.
type Stack struct {
	stack []Model
	name  string
}

// NewStack creates a new [*Stack] with an initial [Model].
// It panics if the initial Model is nil.
func NewStack(initial Model) *Stack {
	if initial == nil {
		panic("initial state cannot be nil")
	}

	return &Stack{
		stack: []Model{initial},
	}
}

// Active returns the Model at the top of the stack.
func (h *Stack) Active() Model {
	if len(h.stack) == 0 {
		return nil
	}
	return h.stack[len(h.stack)-1]
}

// Peek returns the Model at the specified position from the top of the stack.
// If n is out of bounds, it returns nil.
func (h *Stack) Peek(n int) Model {
	if n < 0 || n >= len(h.stack) {
		return nil
	}

	return h.stack[len(h.stack)-1-n]
}

// Depth returns the number of Models in the stack.
func (h *Stack) Depth() int {
	return len(h.stack)
}

func (h *Stack) Name() string { return h.name }

// Init implements [Model].
func (h *Stack) Init() tea.Cmd {
	current := h.Active()
	if initializer, ok := current.(interface{ Init() tea.Cmd }); ok {
		return initializer.Init()
	}
	return nil
}

// Update implements [Model].
func (h *Stack) Update(msg tea.Msg) (Model, tea.Cmd) {
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
		return h, util.Broadcast(ErrStackIsEmpty)
	}

	next, cmd := current.Update(msg)
	h.replace(next)
	return h, cmd
}

// View implements [Model].
func (h *Stack) View() string {
	current := h.Active()
	if current == nil {
		return ""
	}
	return current.View()
}

// push adds a new Model to the top of the stack.
func (h *Stack) push(state Model) {
	if state == nil {
		return
	}

	h.stack = append(h.stack, state)
}

// pop removes the top Model from the stack and returns it.
func (h *Stack) pop() Model {
	if len(h.stack) == 0 {
		return nil
	}

	active := h.Active()
	h.stack = h.stack[:len(h.stack)-1]
	return active
}

// replace the top Model in the stack with a new one.
func (h *Stack) replace(atom Model) {
	if atom == nil {
		return
	}

	if len(h.stack) > 0 {
		h.stack[len(h.stack)-1] = atom
	} else {
		h.stack = append(h.stack, atom)
	}
}

// reset the stack to a single Model state.
// if atom is nil, it resets to the initial Model in the stack.
func (h *Stack) reset(atom Model) {
	if atom == nil {
		h.stack = []Model{h.stack[0]}
		return
	}

	h.stack = []Model{atom}
}
