package polymer

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

// WithLifecycleLogging returns a slice of LensOption that adds logging to the lifecycle events of an Atom.
// It uses the provided logger to log messages for each lifecycle event.
func WithLifecycleLogging(logger *log.Logger) []LensOption {
	return []LensOption{
		WithOnInit(func(active tea.Model, cmd tea.Cmd) {
			logger.Print(formatLog(
				"OnInit",
				fmt.Sprintf("for %s -> %T", formatModel(active), cmd)))
		}),
		WithBeforeUpdate(func(active tea.Model, msg tea.Msg) {
			logger.Print(formatLog(
				"BeforeUpdate",
				fmt.Sprintf("for %s with message <%T> %+v", formatModel(active), msg, msg)))
		}),
		WithAfterUpdate(func(active tea.Model, cmd tea.Cmd) {
			logger.Print(formatLog(
				"AfterUpdate",
				fmt.Sprintf("for %s with command <%T> %+v", formatModel(active), cmd, cmd)))
		}),
		WithOnView(func(active tea.Model, view string) {
			logger.Print(formatLog(
				"OnView",
				fmt.Sprintf("for %s with view not empty <%v>", formatModel(active), view != "")))
		}),
		WithOnError(func(active tea.Model, err error) {
			logger.Print(formatLog(
				"OnError",
				fmt.Sprintf("for %s with error: %v", formatModel(active), err)))
		}),
		WithOnTrace(func(active tea.Model, level trace.Level, msg string) {
			logger.Print(formatLog(
				"OnNotify",
				fmt.Sprintf("[%s] for %s '%s'", level, formatModel(active), msg)))
		}),
	}
}

func formatLog(prefix, msg string) string {
	return fmt.Sprintf("|%-12s| %s", prefix, msg)
}

func formatModel(model tea.Model) string {
	if model == nil {
		return "<nil>"
	}
	if atom, ok := model.(Atomic); ok {
		return fmt.Sprintf("%T (%s)[%d]", model, atom.Name(), atom.Id())
	}
	return fmt.Sprintf("%T", model)
}
