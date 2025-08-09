// Package router provides generic routing utilities for Bubble Tea and [poly.Atomic] models.
// It enables composition and message routing between multiple [Routable] components,
// allowing you to switch active models and delegate updates and rendering.
//
// Main types:
//   - Routable: an interface for Bubble Tea and [poly.Atomic] models that can update and render themselves.
//   - Router: routes messages between two Routable models.
//   - Router3: routes messages between three Routable models.
//   - Slot: identifies which Routable is currently active.
//
// Slot access and mutation helpers:
//   - GetSlotAsT, GetSlotAsU, GetSlotAsV: return pointers to the underlying Routable slot (T, U, or V).
//   - ApplyT, ApplyU, ApplyV: apply a mutation function to the underlying Routable slot.
//
// Example usage:
//
//	r := NewRouter(modelA, modelB, SlotT)
//	r.SetTarget(SlotU)
//	r, cmd := r.Route(msg)
//	view := r.Render()
//	r.ApplyT(func(m *ModelA) { m.State = "updated" })
//	ptr := r.GetSlotAsT()
//
// Also see [menu.Menu] for a practical example of using [Router] with a list model.
package router
