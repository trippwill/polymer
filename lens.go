package polymer

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Lens provides lifecycle hooks for Atoms.
type Lens struct {
	Atom
	OnInit       OnInit       // Called when an Atom is initialized.
	BeforeUpdate BeforeUpdate // Called before an Atom is updated.
	AfterUpdate  AfterUpdate  // Called after an Atom is updated.
	OnView       OnView       // Called when an Atom is rendered.
	OnError      OnError      // Called when an error occurs in an Atom.
	OnNotify     OnNotify     // Called when an Atom sends a notfication.
}

// LensOption configures a Lens.
type LensOption func(*Lens)

// Lifecycle hooks for Atoms.
type (
	OnInit       func(active Atom, cmd tea.Cmd)
	BeforeUpdate func(active Atom, msg tea.Msg)
	AfterUpdate  func(active Atom, cmd tea.Cmd)
	OnView       func(active Atom, rendered string)
	OnError      func(active Atom, err error)
	OnNotify     func(active Atom, level NotificationLevel, msg string)
)

type NotificationLevel int

const (
	TraceLevel NotificationLevel = iota
	DebugLevel
	InfoLevel
	WarnLevel
)

func (n NotificationLevel) String() string {
	switch n {
	case TraceLevel:
		return "Trace"
	case DebugLevel:
		return "Debug"
	case InfoLevel:
		return "Info"
	case WarnLevel:
		return "Warn"
	default:
		return "Unknown"
	}
}

// LensWrap wraps an Atom in a Lens, allowing for lifecycle hooks to be added.
func LensWrap(atom Atom, opts ...LensOption) *Lens {
	l := &Lens{
		Atom: atom,
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

func WithOnNotify(fn OnNotify) LensOption {
	return func(h *Lens) {
		h.OnNotify = fn
	}
}

var _ Atom = Lens{}

// Init implements polymer.Atom.
func (l Lens) Init() tea.Cmd {
	cmd := OptionalInit(l.Atom)
	if l.OnInit != nil {
		l.OnInit(resolve_atom(l.Atom), cmd)
	}

	return cmd
}

// Update implements polymer.Atom.
func (l Lens) Update(msg tea.Msg) (Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case ErrorMsg:
		if l.OnError != nil {
			l.OnError(resolve_atom(l.Atom), msg)
		}
	case NotificationMsg:
		if l.OnNotify != nil {
			l.OnNotify(resolve_atom(l.Atom), msg.level, msg.string)
		}
	}

	if l.BeforeUpdate != nil {
		l.BeforeUpdate(resolve_atom(l.Atom), msg)
	}

	next, cmd := l.Atom.Update(msg)
	if l.AfterUpdate != nil {
		l.AfterUpdate(resolve_atom(next), cmd)
	}

	l.Atom = next
	return l, cmd
}

// View implements polymer.Atom.
func (l Lens) View() string {
	rendered := l.Atom.View()
	if l.OnView != nil {
		l.OnView(resolve_atom(l.Atom), rendered)
	}

	return rendered
}

// resolve_atom recursively resolves an Atom through decorators and lenses.
func resolve_atom(atom Atom) Atom {
	switch a := atom.(type) {
	case *Chain:
		return resolve_atom(a.Active())
	case AtomDecorator:
		return resolve_atom(a.Atom)
	case *Lens:
		return resolve_atom(a.Atom)
	default:
		return a
	}
}
