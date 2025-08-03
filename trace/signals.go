package trace

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/util"
)

// TraceMsg is a tracing message and level.
type TraceMsg struct {
	Msg   string
	Level Level
}

// TraceTrace sends a trace message at the Trace level.
func TraceTrace(msg string) tea.Cmd {
	return util.Broadcast(TraceMsg{Msg: msg, Level: LevelTrace})
}

// TraceDebug sends a debug message at the Debug level.
func TraceDebug(msg string) tea.Cmd {
	return util.Broadcast(TraceMsg{Msg: msg, Level: LevelDebug})
}

// TraceInfo sends an info message at the Info level.
func TraceInfo(msg string) tea.Cmd {
	return util.Broadcast(TraceMsg{Msg: msg, Level: LevelInfo})
}

// TraceWarn sends a warning message at the Warn level.
func TraceWarn(msg string) tea.Cmd {
	return util.Broadcast(TraceMsg{Msg: msg, Level: LevelWarn})
}
