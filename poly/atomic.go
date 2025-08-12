package poly

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/router/multi"
)

// Atomic components are the building blocks of a Polymer application.
// An atomic application is stateful over a single context type X.
// They are integrated into the Bubble Tea model using [Atom]s.
type Atomic[X any] interface {
	Update(msg tea.Msg) (Atomic[X], tea.Cmd)
	View() string
}

type ContextAware[X any] interface {
	SetContext(context X)
}

// Initializer is for components that need to perform initialization.
type Initializer interface {
	Init() tea.Cmd
}

// Identifier is for components that have a unique identifier.
type Identifier interface {
	Id() string
}

// Atom is an adapter for Atomic components integrating them into the Bubble Tea model.
type Atom[X any] struct {
	Model Atomic[X]
}

// ContextMsg carries context data.
type ContextMsg[X any] struct {
	Context X
}

// NewAtom creates a new Atom for the given Atomic model.
func NewAtom[X any](model Atomic[X]) Atom[X] {
	return Atom[X]{
		Model: model,
	}
}

// OptionalInit initializes an [Atomic] model if it implements the [Initializer] interface.
func OptionalInit[X any](model Atomic[X]) tea.Cmd {
	if initModel, ok := model.(Initializer); ok {
		return initModel.Init()
	}
	return nil
}

var _ Identifier = (*Atom[any])(nil)

// Id implements [Identifier].
func (t Atom[X]) Id() string {
	if id, ok := t.Model.(Identifier); ok {
		return id.Id()
	}

	return reflect.TypeOf(t.Model).String()
}

var _ tea.Model = (*Atom[any])(nil)

// Init implements [tea.Model]
func (t Atom[X]) Init() tea.Cmd {
	switch m := t.Model.(type) {
	case Initializer:
		return m.Init()
	default:
		return nil
	}
}

// Update implements [tea.Model].
// [ContextMsg] passes context data to the model.
func (t Atom[X]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if t.Model == nil {
		return nil, nil
	}

	if ctxMsg, ok := msg.(ContextMsg[X]); ok {
		if contextual, ok := t.Model.(ContextAware[X]); ok {
			contextual.SetContext(ctxMsg.Context)
		}
	}

	next, cmd := t.Model.Update(msg)
	if next == nil {
		return nil, cmd
	}

	t.Model = next
	return t, cmd
}

// View implements [tea.Model].
func (t Atom[X]) View() string {
	if t.Model == nil {
		return ""
	}

	return t.Model.View()
}

type AtomicRouter[X any] struct {
	multi.Router[Atomic[X], Atomic[X]]
}

func NewAtomicRouter[X any](initial Atomic[X], fallback Atomic[X], target multi.Slot) AtomicRouter[X] {
	rtr := multi.NewRouter(initial, fallback, target)
	return AtomicRouter[X]{
		Router: rtr,
	}
}

func (r AtomicRouter[T]) Route(msg tea.Msg) (AtomicRouter[T], tea.Cmd) {
	next, cmd := r.Router.Route(msg)
	return AtomicRouter[T]{Router: next}, cmd
}
