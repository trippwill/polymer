package trace

// Level is the severity level of a trace message.
//
//go:generate stringer -type=Level -trimprefix=Level
type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
)
