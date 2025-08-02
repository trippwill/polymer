package trace

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/atom"
)

// WithLifecycleLogging returns a slice of LensOption that adds logging to the lifecycle events of an Atom.
// It uses the provided logger to log messages for each lifecycle event.
func WithLifecycleLogging(logger *log.Logger, minLevel TraceLevel) []LensOption {
	return []LensOption{
		WithOnInit(func(active atom.Model, cmd tea.Cmd) {
			logger.Print(formatLog("OnInit", fmt.Sprintf("for %T -> %T", active, cmd)))
		}),
		WithBeforeUpdate(func(active atom.Model, msg tea.Msg) {
			logger.Print(formatLog("BeforeUpdate", fmt.Sprintf("for %T with message <%T> %+v", active, msg, msg)))
		}),
		WithAfterUpdate(func(active atom.Model, cmd tea.Cmd) {
			logger.Print(formatLog("AfterUpdate", fmt.Sprintf("for %T with command <%T> %+v", active, cmd, cmd)))
		}),
		WithOnView(func(active atom.Model, view string) {
			logger.Print(formatLog("OnView", fmt.Sprintf("for %T with view not empty <%v>", active, view != "")))
		}),
		WithOnError(func(active atom.Model, err error) {
			logger.Print(formatLog("OnError", fmt.Sprintf("for %T with error: %v", active, err)))
		}),
		WithOnTrace(func(active atom.Model, level TraceLevel, msg string) {
			if level >= minLevel {
				logger.Print(formatLog("OnNotify", fmt.Sprintf("[%s] for %T '%s'", level, active, msg)))
			}
		}),
	}
}

func formatLog(prefix, msg string) string {
	return fmt.Sprintf("|%-12s| %s", prefix, msg)
}
