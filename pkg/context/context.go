package context

import (
	"context"

	"github.com/google/uuid"
)

type key int

const (
	requestIDKey key = 0
)

// WithRequestID returns a new Context that carries the request ID.
func WithRequestID(parent context.Context) context.Context {
	return context.WithValue(parent, requestIDKey, uuid.New().String())
}

// RequestID returns the request ID carried by the current Context.
func RequestID(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)[:8]
}
