package multi

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atoms"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// Router holds two router.Routable types and routes messages to one.
type Router[T router.Routable[T], U router.Routable[U]] struct {
	slotT  router.Routable[T]
	slotU  router.Routable[U]
	target Slot
	log    trace.Tracer
}

// NewRouter returns a Routed with the given router.Routables and initial target.
func NewRouter[T router.Routable[T], U router.Routable[U]](
	slotT router.Routable[T],
	slotU router.Routable[U],
	initialTarget Slot,
) Router[T, U] {
	return Router[T, U]{
		slotT:  slotT,
		slotU:  slotU,
		target: initialTarget,
		log:    trace.NewTracer(trace.CategoryRouter),
	}
}

// SetContext sets the context for both SlotT and SlotU if they implement atoms.ContextAware.
func SetContext[T router.Routable[T], U router.Routable[U], X any](r *Router[T, U], ctx X) {
	if r == nil {
		return
	}

	if contextAware, ok := r.slotT.(atoms.ContextAware[X]); ok {
		contextAware.SetContext(ctx)
	}

	if contextAware, ok := r.slotU.(atoms.ContextAware[X]); ok {
		contextAware.SetContext(ctx)
	}
}

// SetTarget sets the active slot.
func (r *Router[T, U]) SetTarget(slot Slot) {
	r.log.Trace("Setting target slot to %s T:%T U:%T", slot, r.slotT, r.slotU)
	r.target = slot
}

// GetTarget returns the current target slot.
func (r Router[T, U]) Target() Slot {
	return r.target
}

// GetT returns the SlotT router.Routable.
func (r Router[T, U]) GetT() router.Routable[T] { return r.slotT }

// GetU returns the SlotU router.Routable.
func (r Router[T, U]) GetU() router.Routable[U] { return r.slotU }

// SetT sets the SlotT router.Routable to the provided value.
func (r Router[T, U]) SetT(value router.Routable[T]) Router[T, U] {
	r.log.Trace("Setting SlotT to %T", value)
	r.slotT = value
	return r
}

// SetU sets the SlotU router.Routable to the provided value.
func (r Router[T, U]) SetU(value router.Routable[U]) Router[T, U] {
	r.log.Trace("Setting SlotU to %T", value)
	r.slotU = value
	return r
}

func (r Router[T, U]) IsSet(slot Slot) bool {
	switch slot {
	case SlotT:
		return r.slotT != nil
	case SlotU:
		return r.slotU != nil
	default:
		r.log.Debug("Unknown slot: %v", slot)
		return false
	}
}

// ConfigureT applies a function to the SlotT router.Routable.
func (r Router[T, U]) ConfigureT(fn func(slot *T)) Router[T, U] {
	r.log.Trace("Applying function to SlotT of type %T", r.slotT)
	if t, ok := r.slotT.(T); ok {
		fn(&t)
		r.slotT = t
	} else {
		r.log.Debug("SlotT is nil, applying function to nil")
		var t T
		fn(&t)
		r.slotT = t
	}
	return r
}

// ConfigureU applies a function to the SlotU router.Routable.
func (r Router[T, U]) ConfigureU(fn func(slot *U)) Router[T, U] {
	r.log.Trace("Applying function to SlotU of type %T", r.slotU)
	if u, ok := r.slotU.(U); ok {
		fn(&u)
		r.slotU = u
	} else {
		r.log.Debug("SlotU is nil, applying function to nil")
		var u U
		fn(&u)
		r.slotU = u
	}
	return r
}

// Route sends the message to the active router.Routable.
func (r Router[T, U]) Route(msg tea.Msg) (Router[T, U], tea.Cmd) {
	switch r.target {
	case SlotT:
		r.log.Trace("Routing message to SlotT of type %T", r.slotT)
		if r.slotT != nil {
			var cmd tea.Cmd
			r.slotT, cmd = r.slotT.Update(msg)
			return r, cmd
		}
	case SlotU:
		r.log.Trace("Routing message to SlotU of type %T", r.slotU)
		if r.slotU != nil {
			var cmd tea.Cmd
			r.slotU, cmd = r.slotU.Update(msg)
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
func (r Router[T, U]) Render() string {
	switch r.target {
	case SlotT:
		r.log.Trace("Rendering SlotT of type %T", r.slotT)
		if r.slotT != nil {
			return r.slotT.View()
		}
	case SlotU:
		r.log.Trace("Rendering SlotU of type %T", r.slotU)
		if r.slotU != nil {
			return r.slotU.View()
		}
	}

	r.log.Trace("No active slot to render view")
	return ""
}
