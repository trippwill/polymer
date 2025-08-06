package router

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
	if r.GetTarget() != SlotT {
		t.Errorf("expected initial target to be SlotT, got %v", r.GetTarget())
	}

	// Set new target
	r.SetTarget(SlotU)
	if r.GetTarget() != SlotU {
		t.Errorf("expected target to be SlotU, got %v", r.GetTarget())
	}
}

func TestRouter_RouteAndRender(t *testing.T) {
	r := NewRouter(mockModel{"a"}, mockModel{"b"}, SlotT)

	// Route to SlotT
	r, _ = r.Route("foo")
	if got := r.SlotT.View(); got != "foo" {
		t.Errorf("expected SlotT view 'foo', got '%s'", got)
	}

	r.SetTarget(SlotU)
	r, _ = r.Route("bar")
	if got := r.SlotU.View(); got != "bar" {
		t.Errorf("expected SlotU view 'bar', got '%s'", got)
	}

	if r.Render() != "bar" {
		t.Errorf("expected render 'bar', got '%s'", r.Render())
	}
}

func TestRouter3_GetAndSetTarget(t *testing.T) {
	r := NewRouter3(mockModel{"a"}, mockModel{"b"}, mockModel{"c"}, SlotV)

	// Check initial target
	if r.GetTarget() != SlotV {
		t.Errorf("expected initial target to be SlotV, got %v", r.GetTarget())
	}

	// Set new target
	r.SetTarget(SlotT)
	if r.GetTarget() != SlotT {
		t.Errorf("expected target to be SlotT, got %v", r.GetTarget())
	}
}

func TestRouter3_RouteAndRender(t *testing.T) {
	r := NewRouter3(mockModel{"a"}, mockModel{"b"}, mockModel{"c"}, SlotV)

	// Route to SlotV
	r, _ = r.Route("baz")
	if got := r.SlotV.View(); got != "baz" {
		t.Errorf("expected SlotV view 'baz', got '%s'", got)
	}

	r.SetTarget(SlotT)
	r, _ = r.Route("foo")
	if got := r.SlotT.View(); got != "foo" {
		t.Errorf("expected SlotT view 'foo', got '%s'", got)
	}

	r.SetTarget(SlotU)
	r, _ = r.Route("bar")
	if got := r.SlotU.View(); got != "bar" {
		t.Errorf("expected SlotU view 'bar', got '%s'", got)
	}

	if r.Render() != "bar" {
		t.Errorf("expected render 'bar', got '%s'", r.Render())
	}
}

func TestApplyAndGetSlot(t *testing.T) {
	r := NewRouter(mockModel{"x"}, mockModel{"y"}, SlotT)

	r.ApplyT(func(m *mockModel) { m.val = "changedT" })
	if r.GetSlotAsT().val != "changedT" {
		t.Errorf("ApplyT did not update SlotT")
	}

	r.ApplyU(func(m *mockModel) { m.val = "changedU" })
	if r.GetSlotAsU().val != "changedU" {
		t.Errorf("ApplyU did not update SlotU")
	}

	r3 := NewRouter3(mockModel{"a"}, mockModel{"b"}, mockModel{"c"}, SlotV)
	r3.ApplyV(func(m *mockModel) { m.val = "changedV" })
	if r3.GetSlotAsV().val != "changedV" {
		t.Errorf("ApplyV did not update SlotV")
	}
}
