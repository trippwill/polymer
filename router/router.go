package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Router holds two Routable types and routes messages to one.
type Router[T Routable[T], U Routable[U]] struct {
	SlotT  Routable[T]
	SlotU  Routable[U]
	target Slot
	log    trace.Tracer
}

// NewRouter returns a Routed with the given Routables and initial target.
func NewRouter[T Routable[T], U Routable[U]](
	slotT Routable[T],
	slotU Routable[U],
	initialTarget Slot,
) Router[T, U] {
	return Router[T, U]{
		SlotT:  slotT,
		SlotU:  slotU,
		target: initialTarget,
		log:    trace.NewTracer(trace.CategoryRouter),
	}
}

// SetTarget sets the active slot.
func (r *Router[T, U]) SetTarget(slot Slot) {
	r.log.Trace("Setting target slot to %s", slot)
	r.target = slot
}

// GetTarget returns the current target slot.
func (r Router[T, U]) GetTarget() Slot {
	return r.target
}

// GetSlotAsT returns a pointer to the current SlotT Routable.
func (r Router[T, U]) GetSlotAsT() *T {
	if slot, ok := r.SlotT.(T); ok {
		return &slot
	}
	r.log.Trace("SlotT is not of type T, returning nil")
	return nil
}

// GetSlotAsU returns a pointer to the current SlotU Routable.
func (r Router[T, U]) GetSlotAsU() *U {
	if slot, ok := r.SlotU.(U); ok {
		return &slot
	}
	r.log.Trace("SlotU is not of type U, returning nil")
	return nil
}

// ApplyT applies a function to the SlotT Routable.
func (r *Router[T, U]) ApplyT(fn func(slot *T)) {
	if t, ok := r.SlotT.(T); ok {
		r.log.Trace("Applying function to SlotT of type %T", t)
		fn(&t)
		r.SlotT = t
	} else {
		r.log.Trace("SlotT is not of type T, cannot apply function")
	}
}

// ApplyU applies a function to the SlotU Routable.
func (r *Router[T, U]) ApplyU(fn func(slot *U)) {
	if u, ok := r.SlotU.(U); ok {
		r.log.Trace("Applying function to SlotU of type %T", u)
		fn(&u)
		r.SlotU = u
	} else {
		r.log.Trace("SlotU is not of type U, cannot apply function")
	}
}

// Route sends the message to the active Routable.
func (r Router[T, U]) Route(msg tea.Msg) (Router[T, U], tea.Cmd) {
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
	default:
		r.log.Trace("Unknown slot: %v", r.target)
		return r, util.Broadcast(ErrUnknownSlot)
	}

	r.log.Trace("No active slot to route message: %T", msg)
	return r, nil
}

// Render returns the active Routable's view.
func (r Router[T, U]) Render() string {
	switch r.target {
	case SlotT:
		r.log.Trace("Rendering SlotT view")
		if r.SlotT != nil {
			return r.SlotT.View()
		} else {
			r.log.Trace("SlotT is nil, cannot render view")
		}
	case SlotU:
		r.log.Trace("Rendering SlotU view")
		if r.SlotU != nil {
			return r.SlotU.View()
		} else {
			r.log.Trace("SlotU is nil, cannot render view")
		}
	}

	r.log.Trace("No active slot to render view")
	return ""
}
