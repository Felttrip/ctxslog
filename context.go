package ctxslog

import (
	"context"
	"sync"
)

type ctxKey string

type safeMap struct {
	mu     sync.Mutex
	fields map[ctxKey]bool
}

var loggedFields = safeMap{}

func init() {
	loggedFields.fields = make(map[ctxKey]bool)
}

func WithValue(ctx context.Context, k string, v any) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	if loggedFields.fields == nil {
		loggedFields.fields = make(map[ctxKey]bool)
	}
	loggedFields.mu.Lock()
	loggedFields.fields[ctxKey(k)] = true
	loggedFields.mu.Unlock()
	return context.WithValue(ctx, ctxKey(k), v)
}

func WithValues(ctx context.Context, fields map[string]any) context.Context {
	if ctx == nil {
		panic("cannot create context from nil parent")
	}
	if loggedFields.fields == nil {
		loggedFields.fields = make(map[ctxKey]bool)
	}
	loggedFields.mu.Lock()
	for k, v := range fields {
		loggedFields.fields[ctxKey(k)] = true
		ctx = context.WithValue(ctx, ctxKey(k), v)
	}
	loggedFields.mu.Unlock()
	return ctx
}
