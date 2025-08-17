package router_test

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/trippwill/polymer/router"
)

// Dummy Routable for testing
type testModel struct {
	val int
}

func (m testModel) Update(msg tea.Msg) (testModel, tea.Cmd) {
	if v, ok := msg.(int); ok {
		m.val += v
	}
	return m, nil
}

func (m testModel) View() string { return "val=" + fmt.Sprint(m.val) }

func TestAutoBasic(t *testing.T) {
	prim := testModel{val: 1}
	ovrd := testModel{val: 10}
	a := router.NewAuto(prim, nil)

	assert.Equal(t, router.AutoSlotPrimary, a.Active())
	assert.Equal(t, "val=1", a.Render())

	a = a.Set(router.AutoSlotOverride, ovrd)
	assert.Equal(t, router.AutoSlotOverride, a.Active())
	assert.Equal(t, "val=10", a.Render())

	// Route to override
	a, _ = a.Route(5)
	assert.Equal(t, "val=15", a.Render())

	// Remove override, route to primary
	a = a.Set(router.AutoSlotOverride, nil)
	a, _ = a.Route(2)
	assert.Equal(t, "val=3", a.Render())
}

func TestAutoConfigure(t *testing.T) {
	a := router.NewAuto(testModel{}, nil)
	a = a.Configure(router.AutoSlotPrimary, func(m *testModel) { m.val = 42 })
	assert.Equal(t, "val=42", a.Render())

	a = a.Set(router.AutoSlotOverride, testModel{})
	a = a.Configure(router.AutoSlotOverride, func(m *testModel) { m.val = 99 })
	assert.Equal(t, "val=99", a.Render())
}
