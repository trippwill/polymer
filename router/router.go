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
	r.log.Trace("Setting target slot to %s T:%T U:%T", slot, r.SlotT, r.SlotU)
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
	r.log.Trace("SlotT is not of type T T:%T, returning nil", r.SlotT)
	return nil
}

// GetSlotAsU returns a pointer to the current SlotU Routable.
func (r Router[T, U]) GetSlotAsU() *U {
	if slot, ok := r.SlotU.(U); ok {
		return &slot
	}
	r.log.Trace("SlotU is not of type U U:%T, returning nil", r.SlotU)
	return nil
}

// ApplyT applies a function to the SlotT Routable.
func (r *Router[T, U]) ApplyT(fn func(slot *T)) {
	r.log.Trace("Applying function to SlotT of type %T", r.SlotT)
	if t, ok := r.SlotT.(T); ok {
		fn(&t)
		r.SlotT = t
	} else {
		var t T
		fn(&t)
		r.SlotT = t
	}
}

// ApplyU applies a function to the SlotU Routable.
func (r *Router[T, U]) ApplyU(fn func(slot *U)) {
	r.log.Trace("Applying function to SlotU of type %T", r.SlotU)
	if u, ok := r.SlotU.(U); ok {
		fn(&u)
		r.SlotU = u
	} else {
		r.log.Debug("SlotU is nil, applying function to nil")
		var u U
		fn(&u)
		r.SlotU = u
	}
}

// Route sends the message to the active Routable.
func (r Router[T, U]) Route(msg tea.Msg) (Router[T, U], tea.Cmd) {
	switch r.target {
	case SlotT:
		r.log.Trace("Routing message to SlotT of type %T", r.SlotT)
		if r.SlotT != nil {
			var cmd tea.Cmd
			r.SlotT, cmd = r.SlotT.Update(msg)
			return r, cmd
		}
	case SlotU:
		r.log.Trace("Routing message to SlotU of type %T", r.SlotU)
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
		r.log.Trace("Rendering SlotT of type %T", r.SlotT)
		if r.SlotT != nil {
			return r.SlotT.View()
		}
	case SlotU:
		r.log.Trace("Rendering SlotU of type %T", r.SlotU)
		if r.SlotU != nil {
			return r.SlotU.View()
		}
	}

	r.log.Trace("No active slot to render view")
	return ""
}
