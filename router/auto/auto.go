package auto

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/poly"
	"github.com/trippwill/polymer/router"
	"github.com/trippwill/polymer/trace"
)

// Slot represents the target slot in the Auto router.
//
//go:generate stringer -type=Slot
type Slot int

const (
	SlotInvalid Slot = iota - 1 // Invalid slot, used for error handling
	SlotPrimary
	SlotOverride
)

// Auto is a generic router that switches between two Routable components.
// It displays the primary Routable unless the override Routable is set.
// This is useful for workflows where a main view is temporarily replaced by an alternate view.
type Auto[T router.Routable[T]] struct {
	primary  router.Routable[T]
	override router.Routable[T]
	log      trace.Tracer
}

// NewAuto constructs an Auto router with a primary and optional override Routable.
// The primary Routable is shown unless the override is set.
func NewAuto[T router.Routable[T]](prim, ovrd router.Routable[T]) Auto[T] {
	return Auto[T]{
		primary:  prim,
		override: ovrd,
		log:      trace.NewTracer(trace.CategoryRouter),
	}
}

// Apply applies a function to the active Routable based on the specified slot,
// returning the result of the function.
func Apply[T router.Routable[T], U any](auto Auto[T], slot Slot, fn func(*T) U) U {
	switch slot {
	case SlotPrimary:
		if pt, ok := auto.primary.(T); ok {
			return fn(&pt)
		}
		return fn(nil)
	case SlotOverride:
		if ot, ok := auto.override.(T); ok {
			return fn(&ot)
		}
		return fn(nil)
	default:
		panic("Invalid Slot provided, must be SlotPrimary or SlotOverride")
	}
}

// Active returns a value indicating which Routable is currently active.
func (a Auto[T]) Active() Slot {
	if a.override != nil {
		return SlotOverride
	}
	if a.primary != nil {
		return SlotPrimary
	}

	return SlotInvalid
}

// IsSet checks if the specified slot (primary or override) has a Routable set.
func (a Auto[T]) IsSet(slot Slot) bool {
	switch slot {
	case SlotPrimary:
		return a.primary != nil
	case SlotOverride:
		return a.override != nil
	default:
		a.log.Warn("Invalid Slot provided, must be SlotPrimary or SlotOverride")
		return false
	}
}

// Set assigns a Routable to the specified slot (primary or override).
func (a Auto[T]) Set(slot Slot, routable router.Routable[T]) Auto[T] {
	switch slot {
	case SlotPrimary:
		a.primary = routable
	case SlotOverride:
		a.override = routable
	default:
		a.log.Warn("Invalid Slot provided, must be SlotPrimary or SlotOverride")
	}
	return a
}

// SetIfNil assigns a Routable to the specified slot (primary or override) only if that slot is currently nil.
func (a Auto[T]) SetIfNil(slot Slot, routable router.Routable[T]) Auto[T] {
	switch slot {
	case SlotPrimary:
		if a.primary == nil {
			a.primary = routable
		}
	case SlotOverride:
		if a.override == nil {
			a.override = routable
		}
	default:
		a.log.Warn("Invalid Slot provided, must be SlotPrimary or SlotOverride")
	}
	return a
}

// Get retrieves the Routable from the specified slot (primary or override).
func (a Auto[T]) Get(slot Slot) router.Routable[T] {
	switch slot {
	case SlotPrimary:
		return a.primary
	case SlotOverride:
		return a.override
	default:
		a.log.Warn("Invalid Slot provided, must be SlotPrimary or SlotOverride")
		return nil
	}
}

// Configure applies a configuration function to the Routable in the specified slot (primary or override).
func (a Auto[T]) Configure(slot Slot, fn func(*T)) Auto[T] {
	if fn == nil {
		a.log.Warn("No configuration function provided, skipping configuration")
		return a
	}
	switch slot {
	case SlotPrimary:
		a.primary = configureRoutable(a.primary, fn)
	case SlotOverride:
		a.override = configureRoutable(a.override, fn)
	default:
		a.log.Warn("Invalid Slot provided, must be SlotPrimary or SlotOverride")
	}
	return a
}

func configureRoutable[T router.Routable[T]](r router.Routable[T], fn func(*T)) router.Routable[T] {
	if r == nil {
		t := new(T)
		fn(t)
		return *t
	}
	if val, ok := r.(T); ok {
		fn(&val)
		return val
	}
	panic("Routable is not of type T, cannot configure")
}

// Route sends a message to the active Routable (override if set, otherwise primary).
// Returns the updated Auto and any resulting command.
func (a Auto[T]) Route(msg tea.Msg) (Auto[T], tea.Cmd) {
	a.log.Trace("Routing message %T", msg)

	if a.override != nil {
		a.log.Trace("Routing to override Routable %T", a.override)
		var cmd tea.Cmd
		a.override, cmd = a.override.Update(msg)
		return a, cmd
	}

	if a.primary != nil {
		a.log.Trace("Routing to primary Routable %T", a.primary)
		var cmd tea.Cmd
		a.primary, cmd = a.primary.Update(msg)
		return a, cmd
	}

	a.log.Warn("No active Routable to route message")
	return a, nil
}

// Render returns the view of the active Routable (override if set, otherwise primary).
func (a Auto[T]) Render() string {
	if a.override != nil {
		a.log.Trace("Rendering override Routable %T", a.override)
		return a.override.View()
	}

	if a.primary != nil {
		a.log.Trace("Rendering primary Routable %T", a.primary)
		return a.primary.View()
	}

	a.log.Warn("No active Routable to render")
	return ""
}

func SetContext[T router.Routable[T], X any](auto *Auto[T], ctx X) {
	if auto == nil {
		return
	}

	if contextAware, ok := any(auto.primary).(poly.ContextAware[X]); ok {
		contextAware.SetContext(ctx)
	}

	if contextAware, ok := any(auto.override).(poly.ContextAware[X]); ok {
		contextAware.SetContext(ctx)
	}
}
