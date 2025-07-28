package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
	"github.com/trippwill/polymer/gels/menu"
	"github.com/trippwill/polymer/observability"
)

// QuitAtom is a simple Atom that quits the application immediately.
type QuitAtom struct{}

func (q QuitAtom) Init() tea.Cmd                           { return tea.Quit }
func (q QuitAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) { return q, nil }
func (q QuitAtom) View() string                            { return "Goodbye!\n" }
func (q QuitAtom) Name() string                            { return "Quit" }

// CounterAtom demonstrates a simple counter with enhanced observability
type CounterAtom struct {
	count int
	max   int
}

func NewCounterAtom(max int) *CounterAtom {
	return &CounterAtom{count: 0, max: max}
}

func (c *CounterAtom) Name() string { return "Counter" }

func (c *CounterAtom) Init() tea.Cmd {
	return nil
}

func (c *CounterAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if c.count < c.max {
				c.count++
			}
		case "down", "j":
			if c.count > 0 {
				c.count--
			}
		case "r":
			c.count = 0
		case "esc":
			return c, poly.Pop()
		case "e":
			// Simulate an error for demonstration
			return c, poly.Error(fmt.Errorf("simulated error at count %d", c.count))
		}
	case poly.ErrorMsg:
		// Handle the error (for demonstration)
		return c, nil
	}
	return c, nil
}

func (c *CounterAtom) View() string {
	return fmt.Sprintf(`
Counter: %d / %d

Controls:
↑/k: Increment
↓/j: Decrement  
r: Reset
e: Simulate error
esc: Back to menu

`, c.count, c.max)
}

// SlowAtom demonstrates performance tracing with artificial delays
type SlowAtom struct {
	counter int
}

func NewSlowAtom() *SlowAtom {
	return &SlowAtom{counter: 0}
}

func (s *SlowAtom) Name() string { return "Slow Operations" }

func (s *SlowAtom) Init() tea.Cmd {
	// Simulate slow initialization
	time.Sleep(100 * time.Millisecond)
	return nil
}

func (s *SlowAtom) Update(msg tea.Msg) (poly.Atom, tea.Cmd) {
	// Simulate slow update processing
	time.Sleep(50 * time.Millisecond)
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			// Simulate very slow operation
			time.Sleep(500 * time.Millisecond)
			s.counter++
		case "f":
			// Fast operation
			s.counter++
		case "esc":
			return s, poly.Pop()
		}
	}
	return s, nil
}

func (s *SlowAtom) View() string {
	// Simulate slow rendering
	time.Sleep(25 * time.Millisecond)
	
	return fmt.Sprintf(`
Slow Operations Demo
Counter: %d

Controls:
s: Slow operation (500ms)
f: Fast operation
esc: Back to menu

This atom demonstrates performance tracing
with artificial delays in init, update, and view.

`, s.counter)
}

func main() {
	// Create log file for demonstration
	logFile, err := os.OpenFile("observability_demo.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		os.Exit(1)
	}
	defer logFile.Close()

	// Set up multiple tracers for comprehensive observability
	
	// 1. Memory tracer for testing and analysis
	_, memoryTracer := observability.TestConfig()
	
	// 2. JSON file tracer for structured logging
	jsonLogger := log.New(logFile, "", log.LstdFlags)
	jsonTracer := observability.NewJSONLoggerTracer(jsonLogger, observability.TraceLevelDebug)
	
	// 3. Console tracer for real-time monitoring
	consoleLogger := log.New(os.Stderr, "[TRACE] ", log.LstdFlags)
	consoleTracer := observability.NewLoggerTracer(consoleLogger, observability.TraceLevelInfo)
	
	// Create comprehensive observability configuration
	config := observability.NewBuilder().
		WithLevel(observability.TraceLevelDebug).
		WithTracer(memoryTracer).
		WithTracer(jsonTracer).
		WithTracer(consoleTracer).
		WithMetadata("app", "observability-demo").
		WithMetadata("version", "1.0.0").
		WithMetadata("session_id", fmt.Sprintf("session-%d", time.Now().Unix())).
		Build()

	// Create atoms with different observability setups
	
	// Counter with full observability
	counter := observability.WithObservability(
		NewCounterAtom(10),
		config,
	)
	counter.SetTraceMetadata("component", "counter")
	counter.SetTraceMetadata("max_value", 10)
	
	// Slow atom with performance-focused tracing
	perfConfig := observability.NewBuilder().
		WithLevel(observability.TraceLevelDebug).
		WithTracer(memoryTracer).
		WithTracer(consoleTracer).
		WithMetadata("performance_test", true).
		WithMetadata("component", "slow-ops").
		Build()
	
	slowAtom := observability.WithObservability(
		NewSlowAtom(),
		perfConfig,
	)
	
	// Create main menu (with lighter observability)
	menuConfig := observability.NewBuilder().
		WithLevel(observability.TraceLevelInfo).
		WithTracer(consoleTracer).
		WithMetadata("component", "main-menu").
		Build()
	
	mainMenu := observability.WithObservability(
		menu.NewMenu(
			"Observability Demo",
			menu.NewMenuItem(counter, "Counter Demo (Full Tracing)"),
			menu.NewMenuItem(slowAtom, "Performance Demo (Timing Focus)"),
			menu.NewMenuItem(QuitAtom{}, "Exit Application"),
		),
		menuConfig,
	)

	// Wrap the menu in a navigation chain with observability
	nav := observability.WithObservability(
		poly.NewChain(mainMenu),
		config,
	)
	nav.SetTraceMetadata("component", "navigation")

	// Create the host and start the application
	host := poly.NewHost(nav, "ObservabilityDemoApp")
	
	fmt.Println("Starting Observability Demo...")
	fmt.Println("- Watch the console for real-time trace output")
	fmt.Println("- Check 'observability_demo.log' for detailed JSON traces")
	fmt.Println("- Memory tracer captures all events for analysis")
	fmt.Println()
	
	p := tea.NewProgram(host)
	
	// Start a goroutine to periodically report memory tracer statistics
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		for range ticker.C {
			events := memoryTracer.GetEvents()
			if len(events) > 0 {
				fmt.Fprintf(os.Stderr, "[STATS] Total events captured: %d\n", len(events))
				
				// Count by event type
				eventCounts := make(map[string]int)
				for _, event := range events {
					eventCounts[event.EventType]++
				}
				
				for eventType, count := range eventCounts {
					fmt.Fprintf(os.Stderr, "[STATS] %s: %d events\n", eventType, count)
				}
				fmt.Fprintf(os.Stderr, "\n")
			}
		}
	}()
	
	// Custom event logging before start
	nav.TraceCustomEvent(
		"app.start",
		"Application starting with full observability enabled",
		observability.TraceLevelInfo,
		map[string]interface{}{
			"startup_time": time.Now().Format(time.RFC3339),
			"tracers":      []string{"memory", "json_file", "console"},
		},
	)
	
	// Run the application
	if _, err := p.Run(); err != nil {
		fmt.Printf("Application failed: %v\n", err)
		os.Exit(1)
	}
	
	// Log final statistics
	events := memoryTracer.GetEvents()
	fmt.Printf("\n=== Final Statistics ===\n")
	fmt.Printf("Total events captured: %d\n", len(events))
	
	// Group by level
	levelCounts := make(map[observability.TraceLevel]int)
	for _, event := range events {
		levelCounts[event.Level]++
	}
	
	fmt.Printf("Events by level:\n")
	for level, count := range levelCounts {
		fmt.Printf("  %s: %d\n", level.String(), count)
	}
	
	// Show performance events
	updateEvents := memoryTracer.GetEventsByType("atom.update")
	var totalUpdateTime time.Duration
	updateCount := 0
	
	for _, event := range updateEvents {
		if event.Duration != nil {
			totalUpdateTime += *event.Duration
			updateCount++
		}
	}
	
	if updateCount > 0 {
		avgUpdateTime := totalUpdateTime / time.Duration(updateCount)
		fmt.Printf("Update performance:\n")
		fmt.Printf("  Total updates: %d\n", updateCount)
		fmt.Printf("  Average update time: %v\n", avgUpdateTime)
	}
	
	fmt.Printf("\nTrace file: observability_demo.log\n")
}