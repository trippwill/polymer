package host

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/trippwill/polymer/trace"
)

// ContextMsg is a message type that carries a context value.
type ContextMsg[X any] struct {
	Context *X
}

func Context[X any](context X) tea.Cmd {
	return func() tea.Msg {
		return ContextMsg[X]{Context: &context}
	}
}

type ContextHost[X any] struct {
	Host
	ctx *X
}

func NewContextHost[X any](name string, root tea.Model, ctx X) *ContextHost[X] {
	if root == nil {
		panic("root state cannot be nil")
	}

	host := &ContextHost[X]{
		Host: Host{
			name:  name,
			state: root,
			log:   trace.NewTracer(trace.CategoryHost),
		},
		ctx: &ctx,
	}

	return host
}

var _ tea.Model = (*ContextHost[any])(nil)

func (h *ContextHost[X]) Init() tea.Cmd {
	h.log.Info(">>>> Initializing context host: " + h.name)
	return tea.Sequence(
		h.Host.Init(),
		Context(*h.ctx),
	)
}

func (h *ContextHost[X]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	h.log.Trace("ContextHost received message: %T {%v}", msg, msg)
	switch msg := msg.(type) {
	case ContextMsg[X]:
		// empty context is a request for current context
		if (msg == ContextMsg[X]{}) || (msg.Context == nil) {
			msg.Context = h.ctx
		}

		// a msg with a context is a request to update the context
		h.ctx = msg.Context
	}

	var cmd tea.Cmd
	next, cmd := h.Host.Update(msg)
	if host, ok := next.(Host); ok {
		h.Host = host
	} else {
		panic("Host must be of type Host")
	}
	return h, cmd
}
