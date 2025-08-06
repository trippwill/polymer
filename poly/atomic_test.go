package poly_test

import (
	"reflect"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
)

// Mock with Identifier
type IdentComponent struct{}

func (i IdentComponent) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) { return i, nil }
func (i IdentComponent) View() string                                   { return "Ident" }
func (i IdentComponent) SetContext(ctx any) poly.Atomic[any]            { return i }
func (i IdentComponent) Id() string                                     { return "custom-id" }
func (i IdentComponent) Init() tea.Cmd                                  { return tea.Println("init") }

// Mock without Identifier
type PlainComponent struct{}

func (p PlainComponent) Update(msg tea.Msg) (poly.Atomic[any], tea.Cmd) { return p, nil }
func (p PlainComponent) View() string                                   { return "Plain" }
func (p PlainComponent) SetContext(ctx any) poly.Atomic[any]            { return p }

func TestAtomId(t *testing.T) {
	atomIdent := poly.NewAtom(IdentComponent{})
	if got, want := atomIdent.Id(), "custom-id"; got != want {
		t.Errorf("IdentComponent Id() = %q, want %q", got, want)
	}

	atomPlain := poly.NewAtom(PlainComponent{})
	wantType := reflect.TypeOf(atomPlain.Model).String()
	if got := atomPlain.Id(); got != wantType {
		t.Errorf("PlainComponent Id() = %q, want %q", got, wantType)
	}
}

func TestAtomView(t *testing.T) {
	atom := poly.NewAtom(PlainComponent{})
	if got, want := atom.View(), "Plain"; got != want {
		t.Errorf("View() = %q, want %q", got, want)
	}
}

func TestAtomInit(t *testing.T) {
	atom := poly.NewAtom(IdentComponent{})
	cmd := atom.Init()
	if cmd == nil {
		t.Error("Init() should not be nil for Initializer")
	}

	atomPlain := poly.NewAtom(PlainComponent{})
	cmdPlain := atomPlain.Init()
	if cmdPlain != nil {
		t.Error("Init() should be nil for non-Initializer")
	}
}

func TestAtomUpdateAndSetContext(t *testing.T) {
	atom := poly.NewAtom(PlainComponent{})
	ctxMsg := poly.ContextMsg[any]{Context: "ctx"}
	updated, _ := atom.Update(ctxMsg)
	if _, ok := updated.(poly.Atom[any]); !ok {
		t.Error("Update() should return Atom[any]")
	}
}

func TestAtomUpdateWithRegularMsg(t *testing.T) {
	atom := poly.NewAtom(PlainComponent{})
	msg := struct{}{}
	updated, cmd := atom.Update(msg)
	if _, ok := updated.(poly.Atom[any]); !ok {
		t.Error("Update() should return Atom[any]")
	}
	if cmd != nil {
		t.Error("Update() should return nil cmd for PlainComponent")
	}
}

func TestAtomViewNilModel(t *testing.T) {
	var atom poly.Atom[any]
	atom.Model = nil
	if got := atom.View(); got != "" {
		t.Errorf("View() with nil Model = %q, want \"\"", got)
	}
}

func TestAtomUpdateNilModel(t *testing.T) {
	var atom poly.Atom[any]
	atom.Model = nil
	updated, cmd := atom.Update(struct{}{})
	if updated != nil {
		t.Error("Update() with nil Model should return nil model")
	}
	if cmd != nil {
		t.Error("Update() with nil Model should return nil cmd")
	}
}
