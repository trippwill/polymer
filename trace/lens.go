package trace

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

// Lens provides lifecycle hooks for Atoms.
type Lens struct {
	atom.Model
	OnInit       OnInit       // Called when an Atom is initialized.
	BeforeUpdate BeforeUpdate // Called before an Atom is updated.
	AfterUpdate  AfterUpdate  // Called after an Atom is updated.
	OnView       OnView       // Called when an Atom is rendered.
	OnError      OnError      // Called when an error occurs in an Atom.
	OnTrace      OnTrace      // Called when an Atom sends a trace message.
}

// LensOption configures a Lens.
type LensOption func(*Lens)

// Lifecycle hooks for Atoms.
type (
	OnInit       func(active atom.Model, cmd tea.Cmd)
	BeforeUpdate func(active atom.Model, msg tea.Msg)
	AfterUpdate  func(active atom.Model, cmd tea.Cmd)
	OnView       func(active atom.Model, rendered string)
	OnError      func(active atom.Model, err error)
	OnTrace      func(active atom.Model, level TraceLevel, msg string)
)

//go:generate stringer -type=TraceLevel -trimprefix=Level
type TraceLevel int

const (
	LevelTrace TraceLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
)

// NewLens wraps an Atom in a Lens, allowing for lifecycle hooks to be added.
func NewLens(atom atom.Model, opts ...LensOption) *Lens {
	l := &Lens{
		Model: atom,
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

var _ atom.Model = Lens{}

// Init implements atom.Model.
func (l Lens) Init() tea.Cmd {
	cmd := atom.OptionalInit(l.Model)
	if l.OnInit != nil {
		l.OnInit(resolve_atom(l.Model), cmd)
	}

	return cmd
}

// Update implements atom.Model.
func (l Lens) Update(msg tea.Msg) (atom.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		if l.OnError != nil {
			l.OnError(resolve_atom(l.Model), msg)
		}
	case TraceMsg:
		if l.OnTrace != nil {
			l.OnTrace(resolve_atom(l.Model), msg.Level, msg.Msg)
		}
	}

	if l.BeforeUpdate != nil {
		l.BeforeUpdate(resolve_atom(l.Model), msg)
	}

	next, cmd := l.Model.Update(msg)
	if l.AfterUpdate != nil {
		l.AfterUpdate(resolve_atom(next), cmd)
	}

	l.Model = next
	return l, cmd
}

// View implements atom.Model.
func (l Lens) View() string {
	rendered := l.Model.View()
	if l.OnView != nil {
		l.OnView(resolve_atom(l.Model), rendered)
	}

	return rendered
}

// resolve_atom recursively resolves an Atom through decorators and lenses.
func resolve_atom(m_atom atom.Model) atom.Model {
	switch a := m_atom.(type) {
	case *atom.Chain:
		return resolve_atom(a.Active())
	case atom.AtomDecorator:
		return resolve_atom(a.Model)
	case *Lens:
		return resolve_atom(a.Model)
	default:
		return a
	}
}
