package trace

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	mer "github.com/trippwill/polymer"
)

// WithLogging returns a slice of mer.LensOption that adds logging to the lifecycle events of an Atom.
// It uses the provided logger to log messages for each lifecycle event.
func WithLogging(logger *log.Logger) []mer.LensOption {
	return []mer.LensOption{
		mer.WithOnInit(func(active mer.Atom, cmd tea.Cmd) {
			logger.Printf("(OnInit) for %T -> %T", active, cmd)
		}),
		mer.WithBeforeUpdate(func(active mer.Atom, msg tea.Msg) {
			logger.Printf("(BeforeUpdate) for %T with message %T", active, msg)
		}),
		mer.WithAfterUpdate(func(active mer.Atom, cmd tea.Cmd) {
			logger.Printf("(AfterUpdate) for %T with command %T", active, cmd)
		}),
		mer.WithOnView(func(active mer.Atom, view string) {
			logger.Printf("(OnView) for %T", active)
		}),
	}
}
