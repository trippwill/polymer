# Observability Demo

This example demonstrates the comprehensive observability and tracing capabilities of the Polymer framework.

## Features Demonstrated

- **Multi-level tracing** (Debug, Info, Warn, Error)
- **Multiple tracer backends** (Memory, File, Console)
- **Structured logging** (JSON and text formats)
- **Performance monitoring** (timing and duration tracking)
- **Error tracing** (automatic error capture and logging)
- **Custom events** (application-specific tracing)
- **Metadata enrichment** (contextual information)
- **Real-time monitoring** (console output and statistics)

## Running the Demo

```bash
cd examples/observability_demo
go run main.go
```

## What the Demo Shows

### Counter Demo
- Full observability enabled
- Traces all lifecycle events (init, update, view)
- Demonstrates error simulation and tracing
- Shows metadata attachment and context propagation

### Performance Demo  
- Focus on timing and performance metrics
- Artificial delays to demonstrate duration tracking
- Shows how slow operations are captured and measured

### Real-time Monitoring
- Console output shows live trace events
- Statistics printed every 5 seconds
- Final performance summary after exit

## Output Files

- **Console**: Real-time trace output with INFO level and above
- **observability_demo.log**: Detailed JSON trace log with all events
- **Memory tracer**: In-memory storage for programmatic analysis

## Key Concepts

### Trace Levels
- `TraceLevelOff`: Disable tracing
- `TraceLevelError`: Only error events
- `TraceLevelWarn`: Warning and error events
- `TraceLevelInfo`: Info, warning, and error events  
- `TraceLevelDebug`: All events (most verbose)

### Tracer Types
- **NullTracer**: No-op tracer for disabled tracing
- **MemoryTracer**: In-memory storage for testing and analysis
- **LoggerTracer**: Text format logging to any `log.Logger`
- **JSONLoggerTracer**: Structured JSON logging
- **FileTracer**: File-based logging with automatic file management
- **CompositeTracer**: Multiple tracers working together

### Event Types
- `atom.init`: Atom initialization
- `atom.update.start`: Before update processing
- `atom.update`: After update completion
- `atom.view`: View rendering
- `atom.error`: Error conditions
- Custom events: Application-defined events

## Configuration Options

The demo shows several configuration approaches:

### Simple Setup
```go
// Quick development tracing
traced := observability.QuickTrace(myAtom, observability.TraceLevelInfo)

// File-based tracing
traced := observability.QuickTraceFile(myAtom, "app.log", observability.TraceLevelDebug)
```

### Advanced Configuration
```go
config := observability.NewBuilder().
    WithLevel(observability.TraceLevelDebug).
    WithFileTracer("app.log", observability.TraceLevelInfo, true).
    WithStderrTracer(observability.TraceLevelError, false).
    WithMetadata("service", "my-app").
    Build()

traced := observability.WithObservability(myAtom, config)
```

### Global Configuration
```go
observability.SetGlobalConfig(observability.ProdConfig("trace.log"))
traced := observability.Trace(myAtom)
```

## Performance Impact

The observability system is designed to be opt-in and have minimal performance impact:

- **Null tracer**: Zero overhead when tracing is disabled
- **Level checking**: Early exit if events won't be processed
- **Lazy evaluation**: Expensive operations only performed when needed
- **Efficient backends**: Memory and file tracers optimized for performance

## Integration with Existing Code

The new observability system maintains full backward compatibility:

```go
// Old way (still works)
traced := polymer.WithLens(myAtom, trace.WithLogging(logger)...)

// New way with enhanced features
config := observability.DevConfig()
traced := config.Wrap(myAtom)

// Bridge to old lens system
lensOptions := config.ToLensOptions()
traced := polymer.WithLens(myAtom, lensOptions...)
```