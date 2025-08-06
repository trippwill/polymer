package router

import (
	"errors"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// Routable can update and render itself.
// Most tea.Models and bubbles implement this.
type Routable[T any] interface {
	Update(msg tea.Msg) (T, tea.Cmd)
	View() string
}

// Slot identifies a Routable in a Routed struct.
type Slot uint8

const (
	SlotSkip Slot = iota // Skip: no slot handles the message.
	SlotT
	SlotU
	SlotV
)

// ErrUnknownSlot signals an unknown slot.
var ErrUnknownSlot error = errors.New("unknown slot")

// Router holds two Routable types and routes messages to one.
type Router[T Routable[T], U Routable[U]] struct {
	SlotT  Routable[T]
	SlotU  Routable[U]
	target Slot
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
	}
}

// SetTarget sets the active slot.
func (r *Router[T, U]) SetTarget(slot Slot) {
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
	return nil
}

// GetSlotAsU returns a pointer to the current SlotU Routable.
func (r Router[T, U]) GetSlotAsU() *U {
	if slot, ok := r.SlotU.(U); ok {
		return &slot
	}
	return nil
}

// ApplyT applies a function to the SlotT Routable.
func (r *Router[T, U]) ApplyT(fn func(slot *T)) {
	log.SetPrefix("Router.ApplyT: ")
	if t, ok := r.SlotT.(T); ok {
		log.Printf("Applying function to SlotT of type %T", t)
		fn(&t)
		r.SlotT = t
	} else {
		log.Printf("SlotT is not of type T, cannot apply function")
	}
}

// ApplyU applies a function to the SlotU Routable.
func (r *Router[T, U]) ApplyU(fn func(slot *U)) {
	log.SetPrefix("Router.ApplyU: ")
	if u, ok := r.SlotU.(U); ok {
		log.Printf("Applying function to SlotU of type %T", u)
		fn(&u)
		r.SlotU = u
	} else {
		log.Printf("SlotU is not of type U, cannot apply function")
	}
}

// Route sends the message to the active Routable.
func (r Router[T, U]) Route(msg tea.Msg) (Router[T, U], tea.Cmd) {
	log.SetPrefix("poly.Router: ")
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
		log.Printf("Unknown slot: %v", r.target)
		return r, util.Broadcast(ErrUnknownSlot)
	}

	log.Printf("No active slot to route message: %T", msg)
	return r, nil
}

// Render returns the active Routable's view.
func (r Router[T, U]) Render() string {
	switch r.target {
	case SlotT:
		if r.SlotT != nil {
			return r.SlotT.View()
		}
	case SlotU:
		if r.SlotU != nil {
			return r.SlotU.View()
		}
	}
	return ""
}

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
