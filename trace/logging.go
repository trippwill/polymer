package trace

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	poly "github.com/trippwill/polymer"
)

// WithBasicLogging returns a slice of mer.LensOption that adds logging to the lifecycle events of an Atom.
// It uses the provided logger to log messages for each lifecycle event.
func WithBasicLogging(logger *log.Logger, notficationLevel poly.NotificationLevel) []poly.LensOption {
	return []poly.LensOption{
		poly.WithOnInit(func(active poly.Atom, cmd tea.Cmd) {
			logger.Print(formatLog("OnInit", fmt.Sprintf("for %T -> %T", active, cmd)))
		}),
		poly.WithBeforeUpdate(func(active poly.Atom, msg tea.Msg) {
			logger.Print(formatLog("BeforeUpdate", fmt.Sprintf("for %T with message <%T> %+v", active, msg, msg)))
		}),
		poly.WithAfterUpdate(func(active poly.Atom, cmd tea.Cmd) {
			logger.Print(formatLog("AfterUpdate", fmt.Sprintf("for %T with command <%T>", active, cmd)))
		}),
		poly.WithOnView(func(active poly.Atom, view string) {
			logger.Print(formatLog("OnView", fmt.Sprintf("for %T", active)))
		}),
		poly.WithOnError(func(active poly.Atom, err error) {
			logger.Print(formatLog("OnError", fmt.Sprintf("for %T with error: %v", active, err)))
		}),
		poly.WithOnNotify(func(active poly.Atom, level poly.NotificationLevel, msg string) {
			if level >= notficationLevel {
				logger.Print(formatLog("OnNotify", fmt.Sprintf("(%s) for %T '%s'", level, active, msg)))
			}
		}),
	}
}

func formatLog(prefix, msg string) string {
	return fmt.Sprintf("|%-12s| %s", prefix, msg)
}
