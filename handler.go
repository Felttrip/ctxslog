package ctxslog

import (
	"context"
	"log/slog"
)

var _ slog.Handler = Handler{}

type Handler struct {
	handler slog.Handler
}

func NewHandler(handler slog.Handler) slog.Handler {
	return Handler{
		handler: handler,
	}
}

// Enabled implements slog.Handler.
func (h Handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.handler.Enabled(ctx, lvl)
}

// Handle implements slog.Handler.
func (h Handler) Handle(ctx context.Context, r slog.Record) error {
	loggedFields.mu.Lock()
	for k, _ := range loggedFields.fields {
		r.AddAttrs(slog.Any(string(k), ctx.Value(k)))
	}
	loggedFields.mu.Unlock()
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements slog.Handler.
func (h Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return Handler{h.handler.WithAttrs(attrs)}
}

// WithGroup implements slog.Handler.
func (h Handler) WithGroup(name string) slog.Handler {
	return h.handler.WithGroup(name)
}
