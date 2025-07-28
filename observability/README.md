# Polymer Observability Package

A comprehensive, extensible observability and tracing system for Polymer TUI applications that provides opt-in monitoring, debugging, and performance analysis capabilities.

## Overview

The observability package extends Polymer's existing Lens system with modern tracing capabilities:

- **Structured tracing** with correlation IDs and hierarchical spans
- **Multi-backend support** for different output formats and destinations  
- **Performance monitoring** with timing and duration tracking
- **Flexible configuration** with preset and builder patterns
- **Opt-in activation** with zero overhead when disabled
- **Backward compatibility** with existing trace package

## Quick Start

### Simple Development Tracing

```go
import "github.com/trippwill/polymer/observability"

// Quick setup for development
traced := observability.QuickTrace(myAtom, observability.TraceLevelInfo)

// File-based tracing
traced := observability.QuickTraceFile(myAtom, "app.log", observability.TraceLevelDebug)

// JSON file tracing
traced := observability.QuickTraceJSON(myAtom, "trace.json", observability.TraceLevelInfo)
```

### Advanced Configuration

```go
// Multi-tracer setup with different outputs
config := observability.NewBuilder().
    WithLevel(observability.TraceLevelDebug).
    WithFileTracer("app.log", observability.TraceLevelInfo, true).      // JSON file
    WithStderrTracer(observability.TraceLevelError, false).             // Console errors
    WithMemoryTracer(observability.TraceLevelDebug).                    // In-memory analysis
    WithMetadata("service", "my-app").
    WithMetadata("version", "1.0.0").
    Build()

traced := observability.WithObservability(myAtom, config)
```

### Global Configuration

```go
// Set up global configuration
observability.SetGlobalConfig(observability.ProdConfig("trace.log"))

// Use global config for all atoms
traced := observability.Trace(myAtom)

// Or override level for specific atoms
traced := observability.TraceWithLevel(myAtom, observability.TraceLevelDebug)
```

## Trace Levels

The system supports hierarchical trace levels:

- `TraceLevelOff` - Disable all tracing (zero overhead)
- `TraceLevelError` - Only error events
- `TraceLevelWarn` - Warning and error events  
- `TraceLevelInfo` - Informational, warning, and error events
- `TraceLevelDebug` - All events including detailed lifecycle tracing

## Tracer Backends

### NullTracer
No-op tracer for disabled tracing with zero performance overhead.

```go
tracer := observability.NewNullTracer()
```

### MemoryTracer
In-memory storage ideal for testing and programmatic analysis.

```go
tracer := observability.NewMemoryTracer(observability.TraceLevelDebug)

// Access events programmatically
events := tracer.GetEvents()
errorEvents := tracer.GetEventsByLevel(observability.TraceLevelError)
updateEvents := tracer.GetEventsByType("atom.update")
```

### LoggerTracer
Text-based logging to any `log.Logger`.

```go
logger := log.New(os.Stderr, "[TRACE] ", log.LstdFlags)
tracer := observability.NewLoggerTracer(logger, observability.TraceLevelInfo)
```

### JSONLoggerTracer
Structured JSON logging for machine processing.

```go
logger := log.New(file, "", log.LstdFlags)  
tracer := observability.NewJSONLoggerTracer(logger, observability.TraceLevelDebug)
```

### FileTracer
Convenient file-based tracing with automatic file management.

```go
tracer, err := observability.NewFileTracer("app.log", observability.TraceLevelInfo, true)
```

### CompositeTracer
Combine multiple tracers for comprehensive observability.

```go
composite := observability.NewCompositeTracer(
    memoryTracer,
    fileTracer, 
    consoleTracer,
)
```

## Event Types

The system automatically traces standard Atom lifecycle events:

- `atom.init` - Atom initialization
- `atom.update.start` - Before update processing  
- `atom.update` - After update completion with timing
- `atom.view` - View rendering with performance metrics
- `atom.error` - Error conditions and exceptions

Custom events can be added:

```go
traced.TraceCustomEvent(
    "user.action",
    "User performed important action", 
    observability.TraceLevelInfo,
    map[string]interface{}{
        "action_type": "button_click",
        "user_id": "12345",
    },
)
```

## Structured Data

Each trace event includes:

```go
type TraceEvent struct {
    // Correlation
    TraceID      string    `json:"trace_id"`      // Request correlation ID
    SpanID       string    `json:"span_id"`       // Individual operation ID  
    ParentSpanID string    `json:"parent_span_id"` // Parent operation ID
    
    // Event details
    Timestamp    time.Time `json:"timestamp"`     // When event occurred
    EventType    string    `json:"event_type"`    // Type of event
    Level        TraceLevel `json:"level"`        // Severity level
    Message      string    `json:"message"`       // Human readable description
    
    // Context
    AtomName     string    `json:"atom_name"`     // Atom identifier
    AtomType     string    `json:"atom_type"`     // Atom type name
    MessageType  string    `json:"message_type"`  // Bubble Tea message type
    CommandType  string    `json:"command_type"`  // Resulting command type
    
    // Performance
    Duration     *time.Duration `json:"duration"` // Operation duration
    
    // Custom data
    Metadata     map[string]interface{} `json:"metadata"` // Additional context
    Error        error `json:"error"`                     // Error information
}
```

## Preset Configurations

### Development
```go
config := observability.DevConfig()  // Debug level, stderr output
```

### Production  
```go
config := observability.ProdConfig("app.log")  // Info level, JSON file
```

### Testing
```go
config, memTracer := observability.TestConfig()  // Debug level, memory storage
```

### Debug
```go
config := observability.DebugConfig(writer)  // Debug level, custom writer
```

## Performance Considerations

The observability system is designed for minimal performance impact:

- **Zero overhead** when using NullTracer or TraceLevelOff
- **Early level checking** prevents expensive operations when disabled
- **Efficient backends** optimized for high-frequency events
- **Lazy evaluation** of expensive string formatting and reflection
- **Configurable detail levels** to balance insight vs. performance

## Backward Compatibility

The new system maintains full compatibility with the existing trace package:

```go
// Old approach (still works)
import "github.com/trippwill/polymer/trace"
traced := polymer.WithLens(myAtom, trace.WithLogging(logger)...)

// New approach with enhanced features
import "github.com/trippwill/polymer/observability"
traced := observability.QuickTrace(myAtom, observability.TraceLevelInfo)

// Bridge from new to old
config := observability.DevConfig()
lensOptions := config.ToLensOptions()
traced := polymer.WithLens(myAtom, lensOptions...)
```

## Migration Guide

### From Basic Lens Hooks

```go
// Before: Manual lens configuration
traced := polymer.WithLens(myAtom,
    polymer.WithOnInit(func(atom polymer.Atom, cmd tea.Cmd) {
        log.Printf("Init: %s", atom.Name())
    }),
    polymer.WithAfterUpdate(func(atom polymer.Atom, cmd tea.Cmd) {
        log.Printf("Update: %s", atom.Name())  
    }),
)

// After: Comprehensive observability
traced := observability.QuickTrace(myAtom, observability.TraceLevelInfo)
```

### From Custom Logging

```go
// Before: Custom logging setup
logger := log.New(file, "app: ", log.LstdFlags)
traced := polymer.WithLens(myAtom, trace.WithLogging(logger)...)

// After: Enhanced structured logging
traced := observability.QuickTraceJSON(myAtom, "app.log", observability.TraceLevelInfo)
```

## Best Practices

### Production Deployment
- Use `TraceLevelInfo` or `TraceLevelError` in production
- Enable JSON logging for structured analysis
- Include relevant metadata (service, version, environment)
- Set up log rotation for file-based tracers

### Development
- Use `TraceLevelDebug` for detailed insights
- Combine console and memory tracers
- Add custom events for important application logic
- Monitor performance metrics during development

### Testing  
- Use MemoryTracer for programmatic verification
- Test error conditions and edge cases
- Verify trace correlation across components
- Benchmark performance impact

### Debugging
- Increase trace level temporarily for issue investigation
- Use custom events to mark important state changes
- Correlate events using TraceID across components
- Export traces for external analysis tools

## Examples

See `/examples/observability_demo/` for a comprehensive demonstration of all features.

## API Reference

### Core Types
- `TraceLevel` - Hierarchical logging levels
- `TraceEvent` - Structured event data
- `TraceContext` - Correlation and metadata context
- `Tracer` - Backend interface for trace output
- `ObservableLens` - Enhanced atom wrapper with observability

### Configuration
- `Builder` - Fluent configuration interface
- `ObservabilityConfig` - Complete configuration structure
- Preset functions: `DevConfig()`, `ProdConfig()`, `TestConfig()`, `DebugConfig()`

### Convenience Functions
- `QuickTrace()`, `QuickTraceFile()`, `QuickTraceJSON()`, `QuickTraceMemory()`
- `Trace()`, `TraceWithLevel()` - Global configuration usage
- `WithObservability()` - Core wrapping function