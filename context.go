package ctxslog

import (
	"context"
	"sync"
)

type ctxKey string

var fieldsKey ctxKey = "ctxslog_fields"

type safeMap struct {
	mu     sync.Mutex
	fields map[string]any
}

func WithValue(ctx context.Context, k string, v any) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	if sm, ok := ctx.Value(fieldsKey).(*safeMap); ok {
		sm.mu.Lock()
		sm.fields[k] = v
		sm.mu.Unlock()
		return context.WithValue(ctx, fieldsKey, sm)
	}
	sm := &safeMap{fields: map[string]any{k: v}}
	return context.WithValue(ctx, fieldsKey, sm)
}

func WithValues(ctx context.Context, fields map[string]any) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	if sm, ok := ctx.Value(fieldsKey).(*safeMap); ok {
		sm.mu.Lock()
		for k, v := range fields {
			sm.fields[k] = v
		}
		sm.mu.Unlock()
		return context.WithValue(ctx, fieldsKey, sm)
	}
	sm := &safeMap{fields: fields}
	return context.WithValue(ctx, fieldsKey, sm)

}
