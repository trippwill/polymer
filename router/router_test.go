package router

import "testing"

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
