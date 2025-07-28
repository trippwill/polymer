package observability

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// ObservableLens extends the existing Lens system with the new observability features
type ObservableLens struct {
	poly.Atom
	tracer     Tracer
	traceCtx   *TraceContext
	atomTracer *AtomTracer
}

// NewObservableLens creates a new observable lens with tracing capabilities
func NewObservableLens(atom poly.Atom, tracer Tracer, level TraceLevel) *ObservableLens {
	traceCtx := NewTraceContext(level)
	atomTracer := NewAtomTracer(tracer, traceCtx)
	
	return &ObservableLens{
		Atom:       atom,
		tracer:     tracer,
		traceCtx:   traceCtx,
		atomTracer: atomTracer,
	}
}

// Name returns the name of the wrapped atom
func (ol *ObservableLens) Name() string {
	return ol.Atom.Name()
}

// Init traces initialization and delegates to the wrapped atom
func (ol *ObservableLens) Init() tea.Cmd {
	cmd := poly.OptionalInit(ol.Atom)
	
	// Trace the initialization
	ol.atomTracer.TraceInit(
		ol.Atom.Name(),
		getAtomType(ol.Atom),
		cmd,
	)
	
	return cmd
}

// Update traces the update lifecycle and delegates to the wrapped atom
func (ol *ObservableLens) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	startTime := time.Now()
	
	// Trace before update
	if ol.tracer.IsEnabled(TraceLevelDebug) {
		event := &TraceEvent{
			TraceID:      ol.traceCtx.TraceID,
			SpanID:       generateSpanID(),
			ParentSpanID: ol.traceCtx.SpanID,
			Timestamp:    startTime,
			EventType:    "atom.update.start",
			Level:        TraceLevelDebug,
			Message:      "Atom update started",
			AtomName:     ol.Atom.Name(),
			AtomType:     getAtomType(ol.Atom),
			MessageType:  getMessageType(msg),
			Metadata:     copyMetadata(ol.traceCtx.Metadata),
		}
		ol.tracer.Trace(context.Background(), event)
	}
	
	// Perform the update
	next, cmd := ol.Atom.Update(msg)
	duration := time.Since(startTime)
	
	// Trace the update completion
	ol.atomTracer.TraceUpdate(
		ol.Atom.Name(),
		getAtomType(ol.Atom),
		msg,
		cmd,
		duration,
	)
	
	// Update the wrapped atom
	ol.Atom = next
	
	// Handle errors
	if errorMsg, ok := msg.(poly.ErrorMsg); ok {
		ol.atomTracer.TraceError(
			ol.Atom.Name(),
			getAtomType(ol.Atom),
			error(errorMsg),
			"Error message received",
		)
	}
	
	return ol, cmd
}

// View traces the view rendering and delegates to the wrapped atom
func (ol *ObservableLens) View() string {
	startTime := time.Now()
	
	view := ol.Atom.View()
	renderDuration := time.Since(startTime)
	
	// Trace the view rendering
	if ol.tracer.IsEnabled(TraceLevelDebug) {
		event := &TraceEvent{
			TraceID:      ol.traceCtx.TraceID,
			SpanID:       generateSpanID(),
			ParentSpanID: ol.traceCtx.SpanID,
			Timestamp:    time.Now(),
			EventType:    "atom.view",
			Level:        TraceLevelDebug,
			Message:      "Atom view rendered",
			AtomName:     ol.Atom.Name(),
			AtomType:     getAtomType(ol.Atom),
			Duration:     &renderDuration,
			Metadata:     copyMetadata(ol.traceCtx.Metadata),
		}
		
		if event.Metadata == nil {
			event.Metadata = make(map[string]interface{})
		}
		event.Metadata["view_length"] = len(view)
		event.Metadata["render_duration_ns"] = renderDuration.Nanoseconds()
		
		ol.tracer.Trace(context.Background(), event)
	}
	
	return view
}

// SetTraceMetadata adds metadata to the trace context
func (ol *ObservableLens) SetTraceMetadata(key string, value interface{}) {
	ol.traceCtx.SetMetadata(key, value)
}

// GetTraceContext returns the current trace context
func (ol *ObservableLens) GetTraceContext() *TraceContext {
	return ol.traceCtx
}

// TraceCustomEvent allows tracing custom events
func (ol *ObservableLens) TraceCustomEvent(eventType, message string, level TraceLevel, metadata map[string]interface{}) {
	if !ol.tracer.IsEnabled(level) {
		return
	}
	
	event := &TraceEvent{
		TraceID:      ol.traceCtx.TraceID,
		SpanID:       generateSpanID(),
		ParentSpanID: ol.traceCtx.SpanID,
		Timestamp:    time.Now(),
		EventType:    eventType,
		Level:        level,
		Message:      message,
		AtomName:     ol.Atom.Name(),
		AtomType:     getAtomType(ol.Atom),
		Metadata:     metadata,
	}
	
	ol.tracer.Trace(context.Background(), event)
}

// ObservabilityConfig holds configuration for observability setup
type ObservabilityConfig struct {
	Level    TraceLevel
	Tracers  []Tracer
	Metadata map[string]interface{}
}

// WithObservability wraps an atom with observability features
func WithObservability(atom poly.Atom, config ObservabilityConfig) *ObservableLens {
	var tracer Tracer
	
	if len(config.Tracers) == 0 {
		// Default to null tracer if none provided
		tracer = NewNullTracer()
	} else if len(config.Tracers) == 1 {
		tracer = config.Tracers[0]
	} else {
		// Use composite tracer for multiple tracers
		tracer = NewCompositeTracer(config.Tracers...)
	}
	
	ol := NewObservableLens(atom, tracer, config.Level)
	
	// Set initial metadata
	if config.Metadata != nil {
		for k, v := range config.Metadata {
			ol.SetTraceMetadata(k, v)
		}
	}
	
	return ol
}

// ToLensOptions converts observability to traditional lens options for backward compatibility
func (ol *ObservableLens) ToLensOptions() []poly.LensOption {
	return []poly.LensOption{
		poly.WithOnInit(func(active poly.Atom, cmd tea.Cmd) {
			ol.atomTracer.TraceInit(active.Name(), getAtomType(active), cmd)
		}),
		poly.WithBeforeUpdate(func(active poly.Atom, msg tea.Msg) {
			if ol.tracer.IsEnabled(TraceLevelDebug) {
				event := &TraceEvent{
					TraceID:      ol.traceCtx.TraceID,
					SpanID:       generateSpanID(),
					ParentSpanID: ol.traceCtx.SpanID,
					Timestamp:    time.Now(),
					EventType:    "atom.update.start",
					Level:        TraceLevelDebug,
					Message:      "Atom update started",
					AtomName:     active.Name(),
					AtomType:     getAtomType(active),
					MessageType:  getMessageType(msg),
					Metadata:     copyMetadata(ol.traceCtx.Metadata),
				}
				ol.tracer.Trace(context.Background(), event)
			}
		}),
		poly.WithAfterUpdate(func(active poly.Atom, cmd tea.Cmd) {
			ol.atomTracer.TraceUpdate(
				active.Name(),
				getAtomType(active),
				nil, // msg not available in this hook
				cmd,
				0, // duration not available in this hook
			)
		}),
		poly.WithOnView(func(active poly.Atom, view string) {
			ol.atomTracer.TraceView(active.Name(), getAtomType(active), len(view))
		}),
	}
}

// TraceLensOptions creates lens options from a tracer for backward compatibility
func TraceLensOptions(tracer Tracer, level TraceLevel) []poly.LensOption {
	traceCtx := NewTraceContext(level)
	atomTracer := NewAtomTracer(tracer, traceCtx)
	
	return []poly.LensOption{
		poly.WithOnInit(func(active poly.Atom, cmd tea.Cmd) {
			atomTracer.TraceInit(active.Name(), getAtomType(active), cmd)
		}),
		poly.WithBeforeUpdate(func(active poly.Atom, msg tea.Msg) {
			if tracer.IsEnabled(TraceLevelDebug) {
				event := &TraceEvent{
					TraceID:      traceCtx.TraceID,
					SpanID:       generateSpanID(),
					ParentSpanID: traceCtx.SpanID,
					Timestamp:    time.Now(),
					EventType:    "atom.update.start",
					Level:        TraceLevelDebug,
					Message:      "Atom update started",
					AtomName:     active.Name(),
					AtomType:     getAtomType(active),
					MessageType:  getMessageType(msg),
					Metadata:     copyMetadata(traceCtx.Metadata),
				}
				tracer.Trace(context.Background(), event)
			}
		}),
		poly.WithAfterUpdate(func(active poly.Atom, cmd tea.Cmd) {
			atomTracer.TraceUpdate(
				active.Name(),
				getAtomType(active),
				nil, // msg not available in this hook
				cmd,
				0, // duration not available in this hook
			)
		}),
		poly.WithOnView(func(active poly.Atom, view string) {
			atomTracer.TraceView(active.Name(), getAtomType(active), len(view))
		}),
	}
}