package trace

import (
	"log"
	"os"
	"slices"
	"strings"
)

var DefaultLevel = LevelTrace

// Category of trace messages.
type Category string

const (
	CategoryOther  Category = "other"
	CategoryHost   Category = "host"
	CategoryRouter Category = "router"
	CategoryAtom   Category = "atom"
	CategoryMenu   Category = "menu"
	CategoryItem   Category = "item"
)

// MinTraceLevel returns the minimum trace level based on the POLY_TRACE_LEVEL environment variable.
func MinTraceLevel() Level {
	switch strings.ToLower(
		os.Getenv("POLY_TRACE_LEVEL")) {
	case "trace":
		return LevelTrace
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	default:
		return DefaultLevel
	}
}

var (
	cachedPolyTraceRaw  string
	cachedPolyTraceList []string
	cachedPolyTraceAll  bool
)

// IsCategoryEnabled checks if tracing is enabled for a specific category.
// Tracing is enabled if POLY_TRACE contains the category name (case-insensitive, underscores/spaces allowed), or "*" for all.
func IsCategoryEnabled(c Category) bool {
	raw := os.Getenv("POLY_TRACE")
	if raw == "" {
		return false
	}
	rawLower := strings.ToLower(raw)
	if rawLower != cachedPolyTraceRaw {
		cachedPolyTraceRaw = rawLower
		cachedPolyTraceAll = rawLower == "*"
		cachedPolyTraceList = nil
		if !cachedPolyTraceAll {
			list := strings.FieldsFunc(rawLower, func(r rune) bool {
				return r == ',' || r == ' '
			})
			for _, item := range list {
				item = strings.TrimSpace(item)
				item = strings.ReplaceAll(item, " ", "_")
				cachedPolyTraceList = append(cachedPolyTraceList, item)
			}
		}
	}

	if cachedPolyTraceAll {
		return true
	}

	normalized := strings.ToLower(strings.ReplaceAll(string(c), " ", "_"))
	return slices.Contains(cachedPolyTraceList, normalized)
}

// Tracer logs messages for a specific category and level.
type Tracer interface {
	// Trace logs a trace message if tracing is enabled and the level is sufficient.
	Trace(format string, args ...any)

	// Debug logs a debug message if tracing is enabled and the level is sufficient.
	Debug(format string, args ...any)

	// Info logs an info message if tracing is enabled and the level is sufficient.
	Info(format string, args ...any)

	// Warn logs a warning message if tracing is enabled and the level is sufficient.
	Warn(format string, args ...any)
}

// activeTracer traces messages for a specific category and level.
type activeTracer struct {
	category Category
	id       string
	minLevel Level
}

// inactiveTracer is a no-op tracer.
type inactiveTracer struct{}

// NewTracer creates a new Tracer for the given category.
func NewTracer(category Category) Tracer {
	if IsCategoryEnabled(category) {
		return &activeTracer{
			category: category,
			minLevel: MinTraceLevel(),
		}
	} else {
		return &inactiveTracer{}
	}
}

func NewTracerWithId(category Category, id string) Tracer {
	if IsCategoryEnabled(category) {
		return &activeTracer{
			category: category,
			id:       id,
			minLevel: MinTraceLevel(),
		}
	} else {
		return &inactiveTracer{}
	}
}

// Trace logs a trace message if tracing is enabled and the level is sufficient.
func (t *activeTracer) Trace(format string, args ...any) {
	if LevelTrace >= t.minLevel {
		log.Printf("[TRACE] %s{%s}: "+format+"\n", append([]any{t.category, t.id}, args...)...)
	}
}

// Debug logs a debug message if tracing is enabled and the level is sufficient.
func (t *activeTracer) Debug(format string, args ...any) {
	if LevelDebug >= t.minLevel {
		log.Printf("[DEBUG] %s{%s}: "+format+"\n", append([]any{t.category, t.id}, args...)...)
	}
}

// Info logs an info message if tracing is enabled and the level is sufficient.
func (t *activeTracer) Info(format string, args ...any) {
	if LevelInfo >= t.minLevel {
		log.Printf("[INFO] %s{%s}: "+format+"\n", append([]any{t.category, t.id}, args...)...)
	}
}

// Warn logs a warning message if tracing is enabled and the level is sufficient.
func (t *activeTracer) Warn(format string, args ...any) {
	if LevelWarn >= t.minLevel {
		log.Printf("[WARN] %s{%s}: "+format+"\n", append([]any{t.category, t.id}, args...)...)
	}
}

// InactiveTracer implements the Tracer interface but does nothing.
func (t *inactiveTracer) Trace(format string, args ...any) {}
func (t *inactiveTracer) Debug(format string, args ...any) {}
func (t *inactiveTracer) Info(format string, args ...any)  {}
func (t *inactiveTracer) Warn(format string, args ...any)  {}
