package observability

import (
	"io"
	"log"
	"os"

	poly "github.com/trippwill/polymer"
)

// Builder provides a fluent interface for configuring observability
type Builder struct {
	config ObservabilityConfig
}

// NewBuilder creates a new observability builder
func NewBuilder() *Builder {
	return &Builder{
		config: ObservabilityConfig{
			Level:    TraceLevelOff,
			Tracers:  make([]Tracer, 0),
			Metadata: make(map[string]interface{}),
		},
	}
}

// WithLevel sets the trace level
func (b *Builder) WithLevel(level TraceLevel) *Builder {
	b.config.Level = level
	return b
}

// WithTracer adds a tracer to the configuration
func (b *Builder) WithTracer(tracer Tracer) *Builder {
	b.config.Tracers = append(b.config.Tracers, tracer)
	return b
}

// WithLoggerTracer adds a logger tracer
func (b *Builder) WithLoggerTracer(logger *log.Logger, level TraceLevel) *Builder {
	tracer := NewLoggerTracer(logger, level)
	return b.WithTracer(tracer)
}

// WithJSONLoggerTracer adds a JSON logger tracer
func (b *Builder) WithJSONLoggerTracer(logger *log.Logger, level TraceLevel) *Builder {
	tracer := NewJSONLoggerTracer(logger, level)
	return b.WithTracer(tracer)
}

// WithFileTracer adds a file tracer
func (b *Builder) WithFileTracer(filename string, level TraceLevel, useJSON bool) *Builder {
	tracer, err := NewFileTracer(filename, level, useJSON)
	if err != nil {
		// Fallback to null tracer on error
		return b.WithTracer(NewNullTracer())
	}
	return b.WithTracer(tracer)
}

// WithWriterTracer adds a writer tracer
func (b *Builder) WithWriterTracer(writer io.Writer, level TraceLevel, useJSON bool) *Builder {
	tracer := NewWriterTracer(writer, level, useJSON)
	return b.WithTracer(tracer)
}

// WithMemoryTracer adds a memory tracer
func (b *Builder) WithMemoryTracer(level TraceLevel) (*Builder, *MemoryTracer) {
	tracer := NewMemoryTracer(level)
	b.WithTracer(tracer)
	return b, tracer
}

// WithStderrTracer adds a tracer that writes to stderr
func (b *Builder) WithStderrTracer(level TraceLevel, useJSON bool) *Builder {
	logger := log.New(os.Stderr, "[TRACE] ", log.LstdFlags)
	if useJSON {
		return b.WithJSONLoggerTracer(logger, level)
	}
	return b.WithLoggerTracer(logger, level)
}

// WithMetadata adds metadata to the configuration
func (b *Builder) WithMetadata(key string, value interface{}) *Builder {
	if b.config.Metadata == nil {
		b.config.Metadata = make(map[string]interface{})
	}
	b.config.Metadata[key] = value
	return b
}

// Build creates the observability configuration
func (b *Builder) Build() ObservabilityConfig {
	return b.config
}

// Wrap wraps an atom with the built configuration
func (b *Builder) Wrap(atom poly.Atom) *ObservableLens {
	return WithObservability(atom, b.config)
}

// ToLensOptions converts the configuration to lens options for backward compatibility
func (b *Builder) ToLensOptions() []poly.LensOption {
	if len(b.config.Tracers) == 0 {
		return []poly.LensOption{} // No tracers, no options
	}
	
	var tracer Tracer
	if len(b.config.Tracers) == 1 {
		tracer = b.config.Tracers[0]
	} else {
		tracer = NewCompositeTracer(b.config.Tracers...)
	}
	
	return TraceLensOptions(tracer, b.config.Level)
}

// Preset configurations for common use cases

// DevConfig creates a development configuration with stderr logging
func DevConfig() *Builder {
	return NewBuilder().
		WithLevel(TraceLevelDebug).
		WithStderrTracer(TraceLevelDebug, false).
		WithMetadata("environment", "development")
}

// ProdConfig creates a production configuration with file logging
func ProdConfig(logFile string) *Builder {
	return NewBuilder().
		WithLevel(TraceLevelInfo).
		WithFileTracer(logFile, TraceLevelInfo, true).
		WithMetadata("environment", "production")
}

// TestConfig creates a test configuration with memory tracing
func TestConfig() (*Builder, *MemoryTracer) {
	builder := NewBuilder().
		WithLevel(TraceLevelDebug).
		WithMetadata("environment", "test")
	
	return builder.WithMemoryTracer(TraceLevelDebug)
}

// DebugConfig creates a debug configuration with detailed logging
func DebugConfig(writer io.Writer) *Builder {
	return NewBuilder().
		WithLevel(TraceLevelDebug).
		WithWriterTracer(writer, TraceLevelDebug, false).
		WithMetadata("environment", "debug")
}

// Convenience functions for quick setup

// QuickTrace provides simple tracing setup for development
func QuickTrace(atom poly.Atom, level TraceLevel) *ObservableLens {
	return DevConfig().
		WithLevel(level).
		Wrap(atom)
}

// QuickTraceFile provides simple file tracing setup
func QuickTraceFile(atom poly.Atom, filename string, level TraceLevel) *ObservableLens {
	return NewBuilder().
		WithLevel(level).
		WithFileTracer(filename, level, false).
		Wrap(atom)
}

// QuickTraceJSON provides simple JSON file tracing setup
func QuickTraceJSON(atom poly.Atom, filename string, level TraceLevel) *ObservableLens {
	return NewBuilder().
		WithLevel(level).
		WithFileTracer(filename, level, true).
		Wrap(atom)
}

// QuickTraceMemory provides simple memory tracing setup for testing
func QuickTraceMemory(atom poly.Atom, level TraceLevel) (*ObservableLens, *MemoryTracer) {
	builder, memTracer := TestConfig()
	return builder.WithLevel(level).Wrap(atom), memTracer
}

// DisableTrace creates an atom with tracing disabled
func DisableTrace(atom poly.Atom) *ObservableLens {
	return NewBuilder().
		WithLevel(TraceLevelOff).
		Wrap(atom)
}

// Global configuration and factory

var globalConfig *Builder

// SetGlobalConfig sets the global observability configuration
func SetGlobalConfig(config *Builder) {
	globalConfig = config
}

// GetGlobalConfig returns the global observability configuration
func GetGlobalConfig() *Builder {
	if globalConfig == nil {
		// Default to development config
		globalConfig = DevConfig()
	}
	return globalConfig
}

// Trace wraps an atom using the global configuration
func Trace(atom poly.Atom) *ObservableLens {
	return GetGlobalConfig().Wrap(atom)
}

// TraceWithLevel wraps an atom using the global configuration with a specific level
func TraceWithLevel(atom poly.Atom, level TraceLevel) *ObservableLens {
	config := GetGlobalConfig().Build()
	config.Level = level
	return WithObservability(atom, config)
}

// Example usage documentation as constants
const (
	ExampleUsageBasic = `
// Basic usage - development tracing to stderr
traced := observability.QuickTrace(myAtom, observability.TraceLevelInfo)

// File tracing
traced := observability.QuickTraceFile(myAtom, "app.log", observability.TraceLevelDebug)

// JSON file tracing  
traced := observability.QuickTraceJSON(myAtom, "trace.json", observability.TraceLevelInfo)

// Memory tracing for testing
traced, memTracer := observability.QuickTraceMemory(myAtom, observability.TraceLevelDebug)
events := memTracer.GetEvents()
`

	ExampleUsageAdvanced = `
// Advanced configuration
config := observability.NewBuilder().
	WithLevel(observability.TraceLevelInfo).
	WithFileTracer("app.log", observability.TraceLevelInfo, true).
	WithStderrTracer(observability.TraceLevelError, false).
	WithMetadata("service", "my-app").
	WithMetadata("version", "1.0.0").
	Build()

traced := observability.WithObservability(myAtom, config)

// Global configuration
observability.SetGlobalConfig(observability.ProdConfig("trace.log"))
traced := observability.Trace(myAtom)
`

	ExampleUsageBackwardCompatible = `
// Backward compatibility with existing Lens system
tracerConfig := observability.DevConfig()
lensOptions := tracerConfig.ToLensOptions()
traced := polymer.WithLens(myAtom, lensOptions...)
`
)