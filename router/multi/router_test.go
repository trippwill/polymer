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
	r := NewRouter(mockModel{"a"}, mockModel{"b"}, SlotT)

	// Check initial target
	if r.Target() != SlotT {
		t.Errorf("expected initial target to be SlotT, got %v", r.Target())
	}

	// Set new target
	r.SetTarget(SlotU)
	if r.Target() != SlotU {
		t.Errorf("expected target to be SlotU, got %v", r.Target())
	}
}

func TestRouter_RouteAndRender(t *testing.T) {
	r := NewRouter(mockModel{"a"}, mockModel{"b"}, SlotT)

	// Route to SlotT
	r, _ = r.Route("foo")
	if got := r.GetT().View(); got != "foo" {
		t.Errorf("expected SlotT view 'foo', got '%s'", got)
	}

	r.SetTarget(SlotU)
	r, _ = r.Route("bar")
	if got := r.GetU().View(); got != "bar" {
		t.Errorf("expected SlotU view 'bar', got '%s'", got)
	}

	if r.Render() != "bar" {
		t.Errorf("expected render 'bar', got '%s'", r.Render())
	}
}
