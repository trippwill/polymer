package observability

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// TestAtom is a simple atom for testing
type TestAtom struct {
	name    string
	counter int
}

func NewTestAtom(name string) *TestAtom {
	return &TestAtom{name: name, counter: 0}
}

func (t *TestAtom) Name() string {
	return t.name
}

func (t *TestAtom) Init() tea.Cmd {
	return nil
}

func (t *TestAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	t.counter++
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" {
			return t, tea.Quit
		}
	case poly.ErrorMsg:
		return t, nil
	}
	
	return t, nil
}

func (t *TestAtom) View() string {
	return "Test view content"
}

// Test TraceLevel
func TestTraceLevel(t *testing.T) {
	tests := []struct {
		level    TraceLevel
		expected string
	}{
		{TraceLevelOff, "OFF"},
		{TraceLevelError, "ERROR"},
		{TraceLevelWarn, "WARN"},
		{TraceLevelInfo, "INFO"},
		{TraceLevelDebug, "DEBUG"},
		{TraceLevel(999), "UNKNOWN"},
	}
	
	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.level.String())
		}
	}
}

// Test TraceContext
func TestTraceContext(t *testing.T) {
	ctx := NewTraceContext(TraceLevelInfo)
	
	if ctx.TraceID == "" {
		t.Error("TraceID should not be empty")
	}
	
	if ctx.SpanID == "" {
		t.Error("SpanID should not be empty")
	}
	
	if ctx.Level != TraceLevelInfo {
		t.Error("Level should be TraceLevelInfo")
	}
	
	// Test metadata
	ctx.SetMetadata("key", "value")
	if ctx.Metadata["key"] != "value" {
		t.Error("Metadata should be set correctly")
	}
	
	// Test child context
	child := ctx.NewChildContext()
	if child.TraceID != ctx.TraceID {
		t.Error("Child should have same TraceID")
	}
	
	if child.SpanID == ctx.SpanID {
		t.Error("Child should have different SpanID")
	}
	
	if child.ParentSpanID != ctx.SpanID {
		t.Error("Child ParentSpanID should be parent's SpanID")
	}
}

// Test NullTracer
func TestNullTracer(t *testing.T) {
	tracer := NewNullTracer()
	
	if tracer.IsEnabled(TraceLevelError) {
		t.Error("NullTracer should never be enabled")
	}
	
	// Should not panic
	event := &TraceEvent{Level: TraceLevelError}
	tracer.Trace(context.Background(), event)
	
	if err := tracer.Close(); err != nil {
		t.Error("NullTracer.Close() should not return error")
	}
}

// Test MemoryTracer
func TestMemoryTracer(t *testing.T) {
	tracer := NewMemoryTracer(TraceLevelInfo)
	
	if !tracer.IsEnabled(TraceLevelError) {
		t.Error("MemoryTracer should be enabled for error level")
	}
	
	if !tracer.IsEnabled(TraceLevelInfo) {
		t.Error("MemoryTracer should be enabled for info level")
	}
	
	if tracer.IsEnabled(TraceLevelDebug) {
		t.Error("MemoryTracer should not be enabled for debug level")
	}
	
	// Test event storage
	event1 := &TraceEvent{
		EventType: "test.event1",
		Level:     TraceLevelInfo,
		Message:   "Test message 1",
	}
	
	event2 := &TraceEvent{
		EventType: "test.event2",
		Level:     TraceLevelError,
		Message:   "Test message 2",
	}
	
	tracer.Trace(context.Background(), event1)
	tracer.Trace(context.Background(), event2)
	
	events := tracer.GetEvents()
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}
	
	// Test filtering
	infoEvents := tracer.GetEventsByLevel(TraceLevelInfo)
	if len(infoEvents) != 1 {
		t.Errorf("Expected 1 info event, got %d", len(infoEvents))
	}
	
	event1Events := tracer.GetEventsByType("test.event1")
	if len(event1Events) != 1 {
		t.Errorf("Expected 1 event1 event, got %d", len(event1Events))
	}
	
	// Test count and clear
	if tracer.Count() != 2 {
		t.Errorf("Expected count 2, got %d", tracer.Count())
	}
	
	tracer.Clear()
	if tracer.Count() != 0 {
		t.Errorf("Expected count 0 after clear, got %d", tracer.Count())
	}
}

// Test LoggerTracer
func TestLoggerTracer(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	tracer := NewLoggerTracer(logger, TraceLevelInfo)
	
	if !tracer.IsEnabled(TraceLevelInfo) {
		t.Error("LoggerTracer should be enabled for info level")
	}
	
	event := &TraceEvent{
		EventType:   "test.event",
		Level:       TraceLevelInfo,
		Message:     "Test message",
		AtomName:    "TestAtom",
		AtomType:    "TestAtom",
		MessageType: "KeyMsg",
	}
	
	tracer.Trace(context.Background(), event)
	
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Error("Output should contain INFO level")
	}
	
	if !strings.Contains(output, "test.event") {
		t.Error("Output should contain event type")
	}
	
	if !strings.Contains(output, "TestAtom") {
		t.Error("Output should contain atom name")
	}
}

// Test JSONLoggerTracer
func TestJSONLoggerTracer(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	tracer := NewJSONLoggerTracer(logger, TraceLevelInfo)
	
	event := &TraceEvent{
		TraceID:     "test-trace-id",
		SpanID:      "test-span-id",
		EventType:   "test.event",
		Level:       TraceLevelInfo,
		Message:     "Test message",
		AtomName:    "TestAtom",
		AtomType:    "TestAtom",
		Timestamp:   time.Now(),
	}
	
	tracer.Trace(context.Background(), event)
	
	output := buf.String()
	
	// Should be valid JSON
	var parsed TraceEvent
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 1 {
		t.Fatal("Expected at least one line of output")
	}
	
	// Remove timestamp prefix added by logger
	jsonPart := lines[0]
	if idx := strings.Index(jsonPart, "{"); idx >= 0 {
		jsonPart = jsonPart[idx:]
	}
	
	if err := json.Unmarshal([]byte(jsonPart), &parsed); err != nil {
		t.Errorf("Output should be valid JSON: %v", err)
	}
	
	if parsed.TraceID != "test-trace-id" {
		t.Error("Parsed JSON should contain correct TraceID")
	}
}

// Test CompositeTracer
func TestCompositeTracer(t *testing.T) {
	memTracer1 := NewMemoryTracer(TraceLevelInfo)
	memTracer2 := NewMemoryTracer(TraceLevelError)
	
	composite := NewCompositeTracer(memTracer1, memTracer2)
	
	// Should be enabled if any tracer is enabled
	if !composite.IsEnabled(TraceLevelInfo) {
		t.Error("Composite should be enabled for info level")
	}
	
	if !composite.IsEnabled(TraceLevelError) {
		t.Error("Composite should be enabled for error level")
	}
	
	if composite.IsEnabled(TraceLevelDebug) {
		t.Error("Composite should not be enabled for debug level")
	}
	
	// Test tracing to multiple tracers
	event := &TraceEvent{
		EventType: "test.event",
		Level:     TraceLevelError,
		Message:   "Test message",
	}
	
	composite.Trace(context.Background(), event)
	
	// Both tracers should receive the event (error level)
	if memTracer1.Count() != 1 {
		t.Errorf("Expected 1 event in tracer1, got %d", memTracer1.Count())
	}
	
	if memTracer2.Count() != 1 {
		t.Errorf("Expected 1 event in tracer2, got %d", memTracer2.Count())
	}
	
	// Test adding tracer
	memTracer3 := NewMemoryTracer(TraceLevelWarn)
	composite.AddTracer(memTracer3)
	
	warnEvent := &TraceEvent{
		EventType: "test.warn",
		Level:     TraceLevelWarn,
		Message:   "Warning message",
	}
	
	composite.Trace(context.Background(), warnEvent)
	
	// Only tracer1 should receive the warn event (info level includes warn)
	if memTracer1.Count() != 2 {
		t.Errorf("Expected 2 events in tracer1, got %d", memTracer1.Count())
	}
	
	// tracer2 (error level) should not receive warn event
	if memTracer2.Count() != 1 {
		t.Errorf("Expected 1 event in tracer2, got %d", memTracer2.Count())
	}
	
	if memTracer3.Count() != 1 {
		t.Errorf("Expected 1 event in tracer3, got %d", memTracer3.Count())
	}
}

// Test ObservableLens
func TestObservableLens(t *testing.T) {
	memTracer := NewMemoryTracer(TraceLevelDebug)
	atom := NewTestAtom("test-atom")
	
	config := ObservabilityConfig{
		Level:   TraceLevelDebug,
		Tracers: []Tracer{memTracer},
	}
	
	lens := WithObservability(atom, config)
	
	if lens.Name() != "test-atom" {
		t.Error("ObservableLens should delegate name to wrapped atom")
	}
	
	// Test Init tracing
	lens.Init()
	
	events := memTracer.GetEvents()
	if len(events) == 0 {
		t.Error("Init should generate trace event")
	}
	
	initEvents := memTracer.GetEventsByType("atom.init")
	if len(initEvents) != 1 {
		t.Errorf("Expected 1 init event, got %d", len(initEvents))
	}
	
	// Test Update tracing
	memTracer.Clear()
	lens.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	
	updateEvents := memTracer.GetEventsByType("atom.update")
	if len(updateEvents) == 0 {
		t.Error("Update should generate trace events")
	}
	
	// Test View tracing
	memTracer.Clear()
	view := lens.View()
	
	if view != "Test view content" {
		t.Error("ObservableLens should delegate view to wrapped atom")
	}
	
	viewEvents := memTracer.GetEventsByType("atom.view")
	if len(viewEvents) == 0 {
		t.Error("View should generate trace events")
	}
	
	// Test metadata
	lens.SetTraceMetadata("test-key", "test-value")
	traceCtx := lens.GetTraceContext()
	if traceCtx.Metadata["test-key"] != "test-value" {
		t.Error("Metadata should be set correctly")
	}
	
	// Test custom event
	memTracer.Clear()
	lens.TraceCustomEvent("custom.event", "Custom message", TraceLevelInfo, map[string]interface{}{
		"custom": "data",
	})
	
	customEvents := memTracer.GetEventsByType("custom.event")
	if len(customEvents) != 1 {
		t.Errorf("Expected 1 custom event, got %d", len(customEvents))
	}
	
	if customEvents[0].Metadata["custom"] != "data" {
		t.Error("Custom event should contain metadata")
	}
}

// Test Configuration Builder
func TestBuilder(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	
	builder := NewBuilder().
		WithLevel(TraceLevelInfo).
		WithLoggerTracer(logger, TraceLevelInfo).
		WithMetadata("env", "test")
	
	config := builder.Build()
	
	if config.Level != TraceLevelInfo {
		t.Error("Builder should set correct level")
	}
	
	if len(config.Tracers) != 1 {
		t.Error("Builder should add tracer")
	}
	
	if config.Metadata["env"] != "test" {
		t.Error("Builder should set metadata")
	}
	
	// Test wrapping atom
	atom := NewTestAtom("test")
	lens := builder.Wrap(atom)
	
	if lens.Name() != "test" {
		t.Error("Builder.Wrap should create functional lens")
	}
}

// Test FileTracer (basic test without actual file I/O)
func TestFileTracerCreation(t *testing.T) {
	// Create temp file
	tmpFile := "/tmp/test-trace.log"
	defer os.Remove(tmpFile)
	
	tracer, err := NewFileTracer(tmpFile, TraceLevelInfo, false)
	if err != nil {
		t.Errorf("NewFileTracer should not error: %v", err)
	}
	
	if tracer == nil {
		t.Error("NewFileTracer should return tracer")
	}
	
	// Clean up
	tracer.Close()
}

// Test Preset Configurations
func TestPresetConfigurations(t *testing.T) {
	// Test DevConfig
	devBuilder := DevConfig()
	devConfig := devBuilder.Build()
	
	if devConfig.Level != TraceLevelDebug {
		t.Error("DevConfig should use debug level")
	}
	
	if devConfig.Metadata["environment"] != "development" {
		t.Error("DevConfig should set environment metadata")
	}
	
	// Test TestConfig
	testBuilder, memTracer := TestConfig()
	testConfig := testBuilder.Build()
	
	if testConfig.Level != TraceLevelDebug {
		t.Error("TestConfig should use debug level")
	}
	
	if memTracer == nil {
		t.Error("TestConfig should return memory tracer")
	}
}

// Test Utility Functions
func TestUtilityFunctions(t *testing.T) {
	atom := NewTestAtom("test")
	
	// Test getAtomType
	atomType := getAtomType(atom)
	if atomType != "TestAtom" {
		t.Errorf("Expected TestAtom, got %s", atomType)
	}
	
	// Test getMessageType
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	msgType := getMessageType(msg)
	if msgType != "KeyMsg" {
		t.Errorf("Expected KeyMsg, got %s", msgType)
	}
	
	// Test copyMetadata
	original := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}
	
	copy := copyMetadata(original)
	copy["key3"] = "value3"
	
	if original["key3"] != nil {
		t.Error("copyMetadata should create independent copy")
	}
	
	if copy["key1"] != "value1" {
		t.Error("copyMetadata should copy values correctly")
	}
}

// Test Error Handling
func TestErrorHandling(t *testing.T) {
	memTracer := NewMemoryTracer(TraceLevelError)
	atom := NewTestAtom("test-atom")
	
	config := ObservabilityConfig{
		Level:   TraceLevelError,
		Tracers: []Tracer{memTracer},
	}
	
	lens := WithObservability(atom, config)
	
	// Send error message
	errorMsg := poly.ErrorMsg(poly.ErrChainEmpty)
	lens.Update(errorMsg)
	
	// Should trace error
	errorEvents := memTracer.GetEventsByType("atom.error")
	if len(errorEvents) == 0 {
		t.Error("Error message should generate error trace event")
	}
}

// Benchmark tests
func BenchmarkObservableLensUpdate(b *testing.B) {
	memTracer := NewMemoryTracer(TraceLevelDebug)
	atom := NewTestAtom("benchmark-atom")
	
	config := ObservabilityConfig{
		Level:   TraceLevelDebug,
		Tracers: []Tracer{memTracer},
	}
	
	lens := WithObservability(atom, config)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		lens.Update(msg)
	}
}

func BenchmarkNullTracerUpdate(b *testing.B) {
	nullTracer := NewNullTracer()
	atom := NewTestAtom("benchmark-atom")
	
	config := ObservabilityConfig{
		Level:   TraceLevelDebug,
		Tracers: []Tracer{nullTracer},
	}
	
	lens := WithObservability(atom, config)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		lens.Update(msg)
	}
}