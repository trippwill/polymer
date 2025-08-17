package multi

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Router holds two [router.Routable] types and routes messages to one.
type Router[A router.Routable[A], B router.Routable[B]] struct {
	slotA  router.Routable[A]
	slotB  router.Routable[B]
	target Slot
	log    trace.Tracer
}

// New returns a [Router] with the given [router.Routable]s and initial target.
func New[A router.Routable[A], B router.Routable[B]](
	slotA router.Routable[A],
	slotB router.Routable[B],
	initialTarget Slot,
) Router[A, B] {
	return Router[A, B]{
		slotA:  slotA,
		slotB:  slotB,
		target: initialTarget,
		log:    trace.NewTracer(trace.CategoryRouter),
	}
}

// SetTarget sets the active slot.
func (r *Router[A, B]) SetTarget(slot Slot) {
	r.log.Trace("Setting target slot to %s T:%T U:%T", slot, r.slotA, r.slotB)
	r.target = slot
}

// GetTarget returns the current target slot.
func (r Router[A, B]) Target() Slot {
	return r.target
}

// GetA returns the SlotA [router.Routable].
func (r Router[A, B]) GetA() router.Routable[A] { return r.slotA }

// GetB returns the SlotB [router.Routable].
func (r Router[A, B]) GetB() router.Routable[B] { return r.slotB }

// SetA sets the SlotA [router.Routable] to the provided value.
func (r Router[A, B]) SetA(value router.Routable[A]) Router[A, B] {
	r.log.Trace("Setting SlotA to %T", value)
	r.slotA = value
	return r
}

// SetB sets the SlotB [router.Routable] to the provided value.
func (r Router[A, B]) SetB(value router.Routable[B]) Router[A, B] {
	r.log.Trace("Setting SlotB to %T", value)
	r.slotB = value
	return r
}

func (r Router[A, B]) IsSet(slot Slot) bool {
	switch slot {
	case SlotA:
		return r.slotA != nil
	case SlotB:
		return r.slotB != nil
	default:
		r.log.Debug("Unknown slot: %v", slot)
		return false
	}
}

// ConfigureA applies a function to the SlotA router.Routable.
func (r Router[A, B]) ConfigureA(fn func(slot *A)) Router[A, B] {
	r.log.Trace("Applying function to SlotA of type %T", r.slotA)
	if a, ok := r.slotA.(A); ok {
		fn(&a)
		r.slotA = a
	} else {
		r.log.Debug("SlotA is nil, applying function to nil")
		var a A
		fn(&a)
		r.slotA = a
	}
	return r
}

// ConfigureB applies a function to the SlotB router.Routable.
func (r Router[A, B]) ConfigureB(fn func(slot *B)) Router[A, B] {
	r.log.Trace("Applying function to SlotB of type %T", r.slotB)
	if b, ok := r.slotB.(B); ok {
		fn(&b)
		r.slotB = b
	} else {
		r.log.Debug("SlotB is nil, applying function to nil")
		var b B
		fn(&b)
		r.slotB = b
	}
	return r
}

// Route sends the message to the active router.Routable.
func (r Router[A, B]) Route(msg tea.Msg) (Router[A, B], tea.Cmd) {
	switch r.target {
	case SlotA:
		r.log.Trace("Routing message to SlotA of type %T", r.slotA)
		if r.slotA != nil {
			var cmd tea.Cmd
			r.slotA, cmd = r.slotA.Update(msg)
			return r, cmd
		}
	case SlotB:
		r.log.Trace("Routing message to SlotB of type %T", r.slotB)
		if r.slotB != nil {
			var cmd tea.Cmd
			r.slotB, cmd = r.slotB.Update(msg)
			return r, cmd
		}
	default:
		r.log.Trace("Unknown slot: %v", r.target)
		return r, util.Broadcast(router.ErrUnknownSlot)
	}

	r.log.Trace("No active slot to route message: %T", msg)
	return r, nil
}

// Render returns the active router.Routable's view.
func (r Router[A, B]) Render() string {
	switch r.target {
	case SlotA:
		r.log.Trace("Rendering SlotA of type %T", r.slotA)
		if r.slotA != nil {
			return r.slotA.View()
		}
	case SlotB:
		r.log.Trace("Rendering SlotB of type %T", r.slotB)
		if r.slotB != nil {
			return r.slotB.View()
		}
	}

	r.log.Trace("No active slot to render view")
	return ""
}
