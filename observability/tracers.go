package observability

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// NullTracer is a no-op tracer for when tracing is disabled
type NullTracer struct{}

// NewNullTracer creates a new null tracer
func NewNullTracer() *NullTracer {
	return &NullTracer{}
}

// IsEnabled always returns false for null tracer
func (nt *NullTracer) IsEnabled(level TraceLevel) bool {
	return false
}

// Trace does nothing for null tracer
func (nt *NullTracer) Trace(ctx context.Context, event *TraceEvent) {
	// no-op
}

// Close does nothing for null tracer
func (nt *NullTracer) Close() error {
	return nil
}

// LoggerTracer wraps a standard logger for structured trace output
type LoggerTracer struct {
	logger    *log.Logger
	level     TraceLevel
	useJSON   bool
	mu        sync.Mutex
}

// LoggerTracerConfig configures the logger tracer
type LoggerTracerConfig struct {
	Logger   *log.Logger
	Level    TraceLevel
	UseJSON  bool
}

// NewLoggerTracer creates a new logger tracer with default configuration
func NewLoggerTracer(logger *log.Logger, level TraceLevel) *LoggerTracer {
	return &LoggerTracer{
		logger:  logger,
		level:   level,
		useJSON: false,
	}
}

// NewJSONLoggerTracer creates a new logger tracer that outputs JSON
func NewJSONLoggerTracer(logger *log.Logger, level TraceLevel) *LoggerTracer {
	return &LoggerTracer{
		logger:  logger,
		level:   level,
		useJSON: true,
	}
}

// IsEnabled returns true if the level is enabled
func (lt *LoggerTracer) IsEnabled(level TraceLevel) bool {
	return level <= lt.level && level != TraceLevelOff
}

// Trace logs the trace event
func (lt *LoggerTracer) Trace(ctx context.Context, event *TraceEvent) {
	if !lt.IsEnabled(event.Level) {
		return
	}
	
	lt.mu.Lock()
	defer lt.mu.Unlock()
	
	if lt.useJSON {
		lt.logJSON(event)
	} else {
		lt.logText(event)
	}
}

// logJSON logs the event as JSON
func (lt *LoggerTracer) logJSON(event *TraceEvent) {
	data, err := json.Marshal(event)
	if err != nil {
		lt.logger.Printf("Failed to marshal trace event: %v", err)
		return
	}
	lt.logger.Printf("%s", string(data))
}

// logText logs the event as formatted text
func (lt *LoggerTracer) logText(event *TraceEvent) {
	msg := fmt.Sprintf("[%s] %s | %s:%s",
		event.Level.String(),
		event.EventType,
		event.AtomType,
		event.AtomName,
	)
	
	if event.MessageType != "" {
		msg += fmt.Sprintf(" | msg:%s", event.MessageType)
	}
	
	if event.Duration != nil {
		msg += fmt.Sprintf(" | duration:%v", *event.Duration)
	}
	
	if event.Error != nil {
		msg += fmt.Sprintf(" | error:%v", event.Error)
	}
	
	msg += fmt.Sprintf(" | %s", event.Message)
	
	lt.logger.Printf("%s", msg)
}

// Close does nothing for logger tracer
func (lt *LoggerTracer) Close() error {
	return nil
}

// MemoryTracer stores trace events in memory for testing and analysis
type MemoryTracer struct {
	events []TraceEvent
	level  TraceLevel
	mu     sync.RWMutex
}

// NewMemoryTracer creates a new memory tracer
func NewMemoryTracer(level TraceLevel) *MemoryTracer {
	return &MemoryTracer{
		events: make([]TraceEvent, 0),
		level:  level,
	}
}

// IsEnabled returns true if the level is enabled
func (mt *MemoryTracer) IsEnabled(level TraceLevel) bool {
	return level <= mt.level && level != TraceLevelOff
}

// Trace stores the trace event in memory
func (mt *MemoryTracer) Trace(ctx context.Context, event *TraceEvent) {
	if !mt.IsEnabled(event.Level) {
		return
	}
	
	mt.mu.Lock()
	defer mt.mu.Unlock()
	
	// Create a copy to avoid issues with pointer sharing
	eventCopy := *event
	if event.Metadata != nil {
		eventCopy.Metadata = copyMetadata(event.Metadata)
	}
	
	mt.events = append(mt.events, eventCopy)
}

// GetEvents returns a copy of all stored events
func (mt *MemoryTracer) GetEvents() []TraceEvent {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	
	events := make([]TraceEvent, len(mt.events))
	copy(events, mt.events)
	return events
}

// GetEventsByType returns events filtered by event type
func (mt *MemoryTracer) GetEventsByType(eventType string) []TraceEvent {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	
	var filtered []TraceEvent
	for _, event := range mt.events {
		if event.EventType == eventType {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// GetEventsByLevel returns events filtered by level
func (mt *MemoryTracer) GetEventsByLevel(level TraceLevel) []TraceEvent {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	
	var filtered []TraceEvent
	for _, event := range mt.events {
		if event.Level == level {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// Clear removes all stored events
func (mt *MemoryTracer) Clear() {
	mt.mu.Lock()
	defer mt.mu.Unlock()
	
	mt.events = mt.events[:0]
}

// Count returns the number of stored events
func (mt *MemoryTracer) Count() int {
	mt.mu.RLock()
	defer mt.mu.RUnlock()
	
	return len(mt.events)
}

// Close does nothing for memory tracer
func (mt *MemoryTracer) Close() error {
	return nil
}

// FileTracer writes trace events to a file
type FileTracer struct {
	file   *os.File
	logger *LoggerTracer
}

// NewFileTracer creates a new file tracer
func NewFileTracer(filename string, level TraceLevel, useJSON bool) (*FileTracer, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open trace file %s: %w", filename, err)
	}
	
	logger := log.New(file, "", log.LstdFlags)
	var loggerTracer *LoggerTracer
	if useJSON {
		loggerTracer = NewJSONLoggerTracer(logger, level)
	} else {
		loggerTracer = NewLoggerTracer(logger, level)
	}
	
	return &FileTracer{
		file:   file,
		logger: loggerTracer,
	}, nil
}

// IsEnabled delegates to the logger tracer
func (ft *FileTracer) IsEnabled(level TraceLevel) bool {
	return ft.logger.IsEnabled(level)
}

// Trace delegates to the logger tracer
func (ft *FileTracer) Trace(ctx context.Context, event *TraceEvent) {
	ft.logger.Trace(ctx, event)
}

// Close closes the file
func (ft *FileTracer) Close() error {
	if ft.file != nil {
		return ft.file.Close()
	}
	return nil
}

// WriterTracer writes trace events to any io.Writer
type WriterTracer struct {
	writer io.Writer
	logger *LoggerTracer
}

// NewWriterTracer creates a new writer tracer
func NewWriterTracer(writer io.Writer, level TraceLevel, useJSON bool) *WriterTracer {
	logger := log.New(writer, "", log.LstdFlags)
	var loggerTracer *LoggerTracer
	if useJSON {
		loggerTracer = NewJSONLoggerTracer(logger, level)
	} else {
		loggerTracer = NewLoggerTracer(logger, level)
	}
	
	return &WriterTracer{
		writer: writer,
		logger: loggerTracer,
	}
}

// IsEnabled delegates to the logger tracer
func (wt *WriterTracer) IsEnabled(level TraceLevel) bool {
	return wt.logger.IsEnabled(level)
}

// Trace delegates to the logger tracer
func (wt *WriterTracer) Trace(ctx context.Context, event *TraceEvent) {
	wt.logger.Trace(ctx, event)
}

// Close does nothing for writer tracer
func (wt *WriterTracer) Close() error {
	return nil
}