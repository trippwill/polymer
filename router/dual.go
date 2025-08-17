package router

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
	"github.com/trippwill/polymer/util"
)

// DualSlot identifies a [Routable] in a [Routed] or [Routed3] struct.
//
//go:generate stringer -type=DualSlot
type DualSlot uint8

const (
	DualSlotSkip DualSlot = iota // Skip: no slot handles the message.
	DualSlotA
	DualSlotB
)

// Dual holds two [Routable] types and routes messages to one.
type Dual[A Routable[A], B Routable[B]] struct {
	slotA  Routable[A]
	slotB  Routable[B]
	target DualSlot
	log    trace.Tracer
}

// New returns a [Dual] with the given [Routable]s and initial target.
func NewDual[A Routable[A], B Routable[B]](
	slotA Routable[A],
	slotB Routable[B],
	initialTarget DualSlot,
) Dual[A, B] {
	return Dual[A, B]{
		slotA:  slotA,
		slotB:  slotB,
		target: initialTarget,
		log:    trace.NewTracer(trace.CategoryRouter),
	}
}

// SetTarget sets the active slot.
func (r *Dual[A, B]) SetTarget(slot DualSlot) {
	r.log.Trace("Setting target slot to %s T:%T U:%T", slot, r.slotA, r.slotB)
	r.target = slot
}

// GetTarget returns the current target slot.
func (r Dual[A, B]) Target() DualSlot {
	return r.target
}

// GetA returns the SlotA [Routable].
func (r Dual[A, B]) GetA() Routable[A] { return r.slotA }

// GetB returns the SlotB [Routable].
func (r Dual[A, B]) GetB() Routable[B] { return r.slotB }

// SetA sets the SlotA [Routable] to the provided value.
func (r Dual[A, B]) SetA(value Routable[A]) Dual[A, B] {
	r.log.Trace("Setting SlotA to %T", value)
	r.slotA = value
	return r
}

// SetB sets the SlotB [Routable] to the provided value.
func (r Dual[A, B]) SetB(value Routable[B]) Dual[A, B] {
	r.log.Trace("Setting SlotB to %T", value)
	r.slotB = value
	return r
}

func (r Dual[A, B]) IsSet(slot DualSlot) bool {
	switch slot {
	case DualSlotA:
		return r.slotA != nil
	case DualSlotB:
		return r.slotB != nil
	default:
		r.log.Debug("Unknown slot: %v", slot)
		return false
	}
}

// ConfigureA applies a function to the SlotA Routable.
func (r Dual[A, B]) ConfigureA(fn func(slot *A)) Dual[A, B] {
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

// ConfigureB applies a function to the SlotB Routable.
func (r Dual[A, B]) ConfigureB(fn func(slot *B)) Dual[A, B] {
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

// Route sends the message to the active Routable.
func (r Dual[A, B]) Route(msg tea.Msg) (Dual[A, B], tea.Cmd) {
	switch r.target {
	case DualSlotA:
		r.log.Trace("Routing message to SlotA of type %T", r.slotA)
		if r.slotA != nil {
			var cmd tea.Cmd
			r.slotA, cmd = r.slotA.Update(msg)
			return r, cmd
		}
	case DualSlotB:
		r.log.Trace("Routing message to SlotB of type %T", r.slotB)
		if r.slotB != nil {
			var cmd tea.Cmd
			r.slotB, cmd = r.slotB.Update(msg)
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
func (r Dual[A, B]) Render() string {
	switch r.target {
	case DualSlotA:
		r.log.Trace("Rendering SlotA of type %T", r.slotA)
		if r.slotA != nil {
			return r.slotA.View()
		}
	case DualSlotB:
		r.log.Trace("Rendering SlotB of type %T", r.slotB)
		if r.slotB != nil {
			return r.slotB.View()
		}
	}

	r.log.Trace("No active slot to render view")
	return ""
}
