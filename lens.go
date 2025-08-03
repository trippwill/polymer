package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

type Modal interface {
	Atomic
	GetCurrent() Atomic // GetCurrent returns the current active model.
}

// Lens provides lifecycle hooks for [tea.Model].
type Lens struct {
	tea.Model
	Atom
	OnInit       OnInit       // Called when an Atom is initialized.
	BeforeUpdate BeforeUpdate // Called before an Atom is updated.
	AfterUpdate  AfterUpdate  // Called after an Atom is updated.
	OnView       OnView       // Called when an Atom is rendered.
	OnError      OnError      // Called when an error occurs in an Atom.
	OnTrace      OnTrace      // Called when an Atom sends a trace message.
}

// LensOption configures a Lens.
type LensOption func(*Lens)

// Lifecycle hooks for [tea.Model].
type (
	OnInit       func(active tea.Model, cmd tea.Cmd)
	BeforeUpdate func(active tea.Model, msg tea.Msg)
	AfterUpdate  func(active tea.Model, cmd tea.Cmd)
	OnView       func(active tea.Model, rendered string)
	OnError      func(active tea.Model, err error)
	OnTrace      func(active tea.Model, level trace.Level, msg string)
)

// NewLens wraps a [tea.Model] in a Lens, allowing for lifecycle hooks to be added.
func NewLens(model tea.Model, opts ...LensOption) *Lens {
	l := &Lens{
		Model: model,
		Atom:  NewAtom("lens"),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// WithOnInit sets the OnInit hook.
func WithOnInit(fn OnInit) LensOption {
	return func(h *Lens) {
		h.OnInit = fn
	}
}

// WithBeforeUpdate sets the BeforeUpdate hook.
func WithBeforeUpdate(fn BeforeUpdate) LensOption {
	return func(h *Lens) {
		h.BeforeUpdate = fn
	}
}

// WithAfterUpdate sets the AfterUpdate hook.
func WithAfterUpdate(fn AfterUpdate) LensOption {
	return func(h *Lens) {
		h.AfterUpdate = fn
	}
}

// WithOnView sets the OnView hook.
func WithOnView(fn OnView) LensOption {
	return func(h *Lens) {
		h.OnView = fn
	}
}

// WithOnError sets the OnError hook.
func WithOnError(fn OnError) LensOption {
	return func(h *Lens) {
		h.OnError = fn
	}
}

// WithOnTrace sets the OnTrace hook.
func WithOnTrace(fn OnTrace) LensOption {
	return func(h *Lens) {
		h.OnTrace = fn
	}
}

var _ Atomic = Lens{}

func (l Lens) Init() tea.Cmd {
	cmd := l.Model.Init()
	if l.OnInit != nil {
		l.OnInit(resolve(l.Model), cmd)
	}

	return cmd
}

func (l Lens) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		if l.OnError != nil {
			l.OnError(resolve(l.Model), msg)
		}
	case trace.TraceMsg:
		if l.OnTrace != nil {
			l.OnTrace(resolve(l.Model), msg.Level, msg.Msg)
		}
	}

	if l.BeforeUpdate != nil {
		l.BeforeUpdate(resolve(l.Model), msg)
	}

	next, cmd := l.Model.Update(msg)
	if l.AfterUpdate != nil {
		l.AfterUpdate(resolve(next), cmd)
	}

	l.Model = next
	return l, cmd
}

func (l Lens) View() string {
	rendered := ""
	if l.Model != nil {
		rendered = l.Model.View()
	}

	if l.OnView != nil {
		l.OnView(resolve(l.Model), rendered)
	}

	return rendered
}

// resolve recursively resolves a [tea.Model] through [Lens] and [Modal] types.
func resolve(model tea.Model) tea.Model {
	switch a := model.(type) {
	case Modal:
		current := a.GetCurrent()
		if current.Id() != a.Id() {
			// If the current model is not the same as the modal's ID, resolve it.
			return resolve(current)
		}
		return current
	case *Lens:
		return resolve(a.Model)
	case Lens:
		return resolve(a.Model)
	default:
		return a
	}
}
