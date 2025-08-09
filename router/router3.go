package router

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// Router3 holds three Routable types and routes messages to one.
type Router3[T Routable[T], U Routable[U], V Routable[V]] struct {
	Router[T, U]
	SlotV Routable[V]
}

// NewRouter3 returns a Routed3 with the given Routables and initial target.
func NewRouter3[T Routable[T], U Routable[U], V Routable[V]](
	slotT Routable[T],
	slotU Routable[U],
	slotV Routable[V],
	initialTarget Slot,
) Router3[T, U, V] {
	return Router3[T, U, V]{
		Router: NewRouter(slotT, slotU, initialTarget),
		SlotV:  slotV,
	}
}

// GetSlotAsV returns a pointer to the current SlotV Routable.
func (r Router3[T, U, V]) GetSlotAsV() *V {
	if slot, ok := r.SlotV.(V); ok {
		return &slot
	}

	return nil
}

// ApplyV applies a function to the SlotV Routable.
func (r *Router3[T, U, V]) ApplyV(fn func(slot *V)) {
	log.SetPrefix("Router3.ApplyV: ")
	if v, ok := r.SlotV.(V); ok {
		log.Printf("Applying function to SlotV of type %T", v)
		fn(&v)
		r.SlotV = v
	} else {
		log.Printf("SlotV is not of type V, cannot apply function")
	}
}

// Route sends the message to the active Routable.
func (r Router3[T, U, V]) Route(msg tea.Msg) (Router3[T, U, V], tea.Cmd) {
	switch r.target {
	case SlotT:
		if r.SlotT != nil {
			var cmd tea.Cmd
			r.SlotT, cmd = r.SlotT.Update(msg)
			return r, cmd
		}
	case SlotU:
		if r.SlotU != nil {
			var cmd tea.Cmd
			r.SlotU, cmd = r.SlotU.Update(msg)
			return r, cmd
		}
	case SlotV:
		if r.SlotV != nil {
			var cmd tea.Cmd
			r.SlotV, cmd = r.SlotV.Update(msg)
			return r, cmd
		}
	default:
		return r, util.Broadcast(ErrUnknownSlot)
	}
	return r, nil
}

// Render returns the active Routable's view.
func (r Router3[T, U, V]) Render() string {
	switch r.target {
	case SlotT:
		if r.SlotT != nil {
			return r.SlotT.View()
		}
	case SlotU:
		if r.SlotU != nil {
			return r.SlotU.View()
		}
	case SlotV:
		if r.SlotV != nil {
			return r.SlotV.View()
		}
	}
	return ""
}
