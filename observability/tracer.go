package observability

import (
	"context"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// TraceLevel defines the level of tracing detail
type TraceLevel int

const (
	TraceLevelOff TraceLevel = iota
	TraceLevelError
	TraceLevelWarn
	TraceLevelInfo
	TraceLevelDebug
)

// String returns the string representation of the trace level
func (t TraceLevel) String() string {
	switch t {
	case TraceLevelOff:
		return "OFF"
	case TraceLevelError:
		return "ERROR"
	case TraceLevelWarn:
		return "WARN"
	case TraceLevelInfo:
		return "INFO"
	case TraceLevelDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

// TraceEvent represents a structured tracing event
type TraceEvent struct {
	// Core identification
	TraceID     string                 `json:"trace_id"`
	SpanID      string                 `json:"span_id"`
	ParentSpanID string                `json:"parent_span_id,omitempty"`
	
	// Event details
	Timestamp   time.Time              `json:"timestamp"`
	EventType   string                 `json:"event_type"`
	Level       TraceLevel             `json:"level"`
	Message     string                 `json:"message"`
	
	// Context
	AtomName    string                 `json:"atom_name,omitempty"`
	AtomType    string                 `json:"atom_type,omitempty"`
	MessageType string                 `json:"message_type,omitempty"`
	CommandType string                 `json:"command_type,omitempty"`
	
	// Performance
	Duration    *time.Duration         `json:"duration,omitempty"`
	
	// Additional data
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Error       error                  `json:"error,omitempty"`
}

// TraceContext carries tracing state and correlation information
type TraceContext struct {
	TraceID     string
	SpanID      string
	ParentSpanID string
	Level       TraceLevel
	StartTime   time.Time
	Metadata    map[string]interface{}
}

// NewTraceContext creates a new trace context with a unique trace ID
func NewTraceContext(level TraceLevel) *TraceContext {
	return &TraceContext{
		TraceID:   generateTraceID(),
		SpanID:    generateSpanID(),
		Level:     level,
		StartTime: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NewChildContext creates a child span context
func (tc *TraceContext) NewChildContext() *TraceContext {
	return &TraceContext{
		TraceID:      tc.TraceID,
		SpanID:       generateSpanID(),
		ParentSpanID: tc.SpanID,
		Level:        tc.Level,
		StartTime:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}
}

// SetMetadata adds metadata to the trace context
func (tc *TraceContext) SetMetadata(key string, value interface{}) {
	if tc.Metadata == nil {
		tc.Metadata = make(map[string]interface{})
	}
	tc.Metadata[key] = value
}

// Tracer defines the interface for tracing backends
type Tracer interface {
	// IsEnabled checks if tracing is enabled for the given level
	IsEnabled(level TraceLevel) bool
	
	// Trace emits a trace event
	Trace(ctx context.Context, event *TraceEvent)
	
	// Close cleans up the tracer resources
	Close() error
}

// CompositeTracer allows multiple tracers to be used together
type CompositeTracer struct {
	tracers []Tracer
}

// NewCompositeTracer creates a new composite tracer
func NewCompositeTracer(tracers ...Tracer) *CompositeTracer {
	return &CompositeTracer{tracers: tracers}
}

// AddTracer adds a tracer to the composite
func (ct *CompositeTracer) AddTracer(tracer Tracer) {
	ct.tracers = append(ct.tracers, tracer)
}

// IsEnabled returns true if any tracer is enabled for the given level
func (ct *CompositeTracer) IsEnabled(level TraceLevel) bool {
	for _, tracer := range ct.tracers {
		if tracer.IsEnabled(level) {
			return true
		}
	}
	return false
}

// Trace sends the event to all enabled tracers
func (ct *CompositeTracer) Trace(ctx context.Context, event *TraceEvent) {
	for _, tracer := range ct.tracers {
		if tracer.IsEnabled(event.Level) {
			tracer.Trace(ctx, event)
		}
	}
}

// Close closes all tracers
func (ct *CompositeTracer) Close() error {
	var lastErr error
	for _, tracer := range ct.tracers {
		if err := tracer.Close(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// AtomTracer provides convenient methods for tracing Atom lifecycle events
type AtomTracer struct {
	tracer Tracer
	ctx    *TraceContext
}

// NewAtomTracer creates a new atom tracer
func NewAtomTracer(tracer Tracer, ctx *TraceContext) *AtomTracer {
	return &AtomTracer{
		tracer: tracer,
		ctx:    ctx,
	}
}

// TraceInit traces an atom initialization event
func (at *AtomTracer) TraceInit(atomName, atomType string, cmd tea.Cmd) {
	if !at.tracer.IsEnabled(TraceLevelInfo) {
		return
	}
	
	event := &TraceEvent{
		TraceID:     at.ctx.TraceID,
		SpanID:      at.ctx.SpanID,
		ParentSpanID: at.ctx.ParentSpanID,
		Timestamp:   time.Now(),
		EventType:   "atom.init",
		Level:       TraceLevelInfo,
		Message:     "Atom initialized",
		AtomName:    atomName,
		AtomType:    atomType,
		CommandType: getCommandType(cmd),
		Metadata:    copyMetadata(at.ctx.Metadata),
	}
	
	at.tracer.Trace(context.Background(), event)
}

// TraceUpdate traces an atom update event
func (at *AtomTracer) TraceUpdate(atomName, atomType string, msg tea.Msg, cmd tea.Cmd, duration time.Duration) {
	if !at.tracer.IsEnabled(TraceLevelDebug) {
		return
	}
	
	event := &TraceEvent{
		TraceID:     at.ctx.TraceID,
		SpanID:      generateSpanID(), // New span for each update
		ParentSpanID: at.ctx.SpanID,
		Timestamp:   time.Now(),
		EventType:   "atom.update",
		Level:       TraceLevelDebug,
		Message:     "Atom updated",
		AtomName:    atomName,
		AtomType:    atomType,
		MessageType: getMessageType(msg),
		CommandType: getCommandType(cmd),
		Duration:    &duration,
		Metadata:    copyMetadata(at.ctx.Metadata),
	}
	
	at.tracer.Trace(context.Background(), event)
}

// TraceView traces an atom view render event
func (at *AtomTracer) TraceView(atomName, atomType string, viewLength int) {
	if !at.tracer.IsEnabled(TraceLevelDebug) {
		return
	}
	
	event := &TraceEvent{
		TraceID:     at.ctx.TraceID,
		SpanID:      generateSpanID(),
		ParentSpanID: at.ctx.SpanID,
		Timestamp:   time.Now(),
		EventType:   "atom.view",
		Level:       TraceLevelDebug,
		Message:     "Atom view rendered",
		AtomName:    atomName,
		AtomType:    atomType,
		Metadata:    copyMetadata(at.ctx.Metadata),
	}
	
	if event.Metadata == nil {
		event.Metadata = make(map[string]interface{})
	}
	event.Metadata["view_length"] = viewLength
	
	at.tracer.Trace(context.Background(), event)
}

// TraceError traces an error event
func (at *AtomTracer) TraceError(atomName, atomType string, err error, message string) {
	if !at.tracer.IsEnabled(TraceLevelError) {
		return
	}
	
	event := &TraceEvent{
		TraceID:     at.ctx.TraceID,
		SpanID:      generateSpanID(),
		ParentSpanID: at.ctx.SpanID,
		Timestamp:   time.Now(),
		EventType:   "atom.error",
		Level:       TraceLevelError,
		Message:     message,
		AtomName:    atomName,
		AtomType:    atomType,
		Error:       err,
		Metadata:    copyMetadata(at.ctx.Metadata),
	}
	
	at.tracer.Trace(context.Background(), event)
}