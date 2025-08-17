package router_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/router"
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
	r := router.NewDual(mockModel{"a"}, mockModel{"b"}, router.DualSlotA)

	// Check initial target
	if r.Target() != router.DualSlotA {
		t.Errorf("expected initial target to be router.DualSlotT, got %v", r.Target())
	}

	// Set new target
	r.SetTarget(router.DualSlotB)
	if r.Target() != router.DualSlotB {
		t.Errorf("expected target to be router.DualSlotU, got %v", r.Target())
	}
}

func TestRouter_RouteAndRender(t *testing.T) {
	r := router.NewDual(mockModel{"a"}, mockModel{"b"}, router.DualSlotA)

	// Route to router.DualSlotA
	r, _ = r.Route("foo")
	if got := r.GetA().View(); got != "foo" {
		t.Errorf("expected router.DualSlotA view 'foo', got '%s'", got)
	}

	r.SetTarget(router.DualSlotB)
	r, _ = r.Route("bar")
	if got := r.GetB().View(); got != "bar" {
		t.Errorf("expected router.DualSlotB view 'bar', got '%s'", got)
	}

	if r.Render() != "bar" {
		t.Errorf("expected render 'bar', got '%s'", r.Render())
	}
}
