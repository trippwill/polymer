package multi

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type mockModel struct {
	val string
}

func (m mockModel) Update(msg tea.Msg) (mockModel, tea.Cmd) {
	return mockModel{val: msg.(string)}, nil
}

func (m mockModel) View() string {
	return m.val
}

func TestRouter_GetAndSetTarget(t *testing.T) {
	r := New(mockModel{"a"}, mockModel{"b"}, SlotA)

	// Check initial target
	if r.Target() != SlotA {
		t.Errorf("expected initial target to be SlotT, got %v", r.Target())
	}

	// Set new target
	r.SetTarget(SlotB)
	if r.Target() != SlotB {
		t.Errorf("expected target to be SlotU, got %v", r.Target())
	}
}

func TestRouter_RouteAndRender(t *testing.T) {
	r := New(mockModel{"a"}, mockModel{"b"}, SlotA)

	// Route to SlotA
	r, _ = r.Route("foo")
	if got := r.GetA().View(); got != "foo" {
		t.Errorf("expected SlotA view 'foo', got '%s'", got)
	}

	r.SetTarget(SlotB)
	r, _ = r.Route("bar")
	if got := r.GetB().View(); got != "bar" {
		t.Errorf("expected SlotB view 'bar', got '%s'", got)
	}

	if r.Render() != "bar" {
		t.Errorf("expected render 'bar', got '%s'", r.Render())
	}
}
