package polymer

import tea "github.com/charmbracelet/bubbletea"

type (
	OnInit       func(active Atom, cmd tea.Cmd)
	BeforeUpdate func(active Atom, msg tea.Msg)
	AfterUpdate  func(active Atom, cmd tea.Cmd)
	OnView       func(active Atom, rendered string)
)

// Lens provides lifecycle hooks for Atoms.
type Lens struct {
	Atom
	OnInit       OnInit       // Called when an Atom is initialized.
	BeforeUpdate BeforeUpdate // Called before an Atom is updated.
	AfterUpdate  AfterUpdate  // Called after an Atom is updated.
	OnView       OnView       // Called when an Atom is rendered.
}

// LensOption configures a Conduit.
type LensOption func(*Lens)

var _ Atom = Lens{}

func (l Lens) Init() tea.Cmd {
	cmd := OptionalInit(l.Atom)
	if l.OnInit != nil {
		l.OnInit(l.Atom, cmd)
	}

	return cmd
}

func (l Lens) Update(msg tea.Msg) (Atom, tea.Cmd) {
	if l.BeforeUpdate != nil {
		l.BeforeUpdate(l.Atom, msg)
	}

	next, cmd := l.Atom.Update(msg)
	if l.AfterUpdate != nil {
		l.AfterUpdate(next, cmd)
	}

	l.Atom = next
	return l, cmd
}

func (l Lens) View() string {
	rendered := l.Atom.View()
	if l.OnView != nil {
		l.OnView(l.Atom, rendered)
	}

	return rendered
}

func WithLens(atom Atom, opts ...LensOption) *Lens {
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

// WithBeforeUpdate sets the OnMsg hook.
func WithBeforeUpdate(fn BeforeUpdate) LensOption {
	return func(h *Lens) {
		h.BeforeUpdate = fn
	}
}

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
