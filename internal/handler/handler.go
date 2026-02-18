package handler

import (
	"context"

	"github.com/octopunkio/taskflow-worker/internal/job"
)

// Handler processes jobs of a specific type
type Handler interface {
	// Handle processes a job and returns an error if processing fails
	Handle(ctx context.Context, j *job.Job) error
}

// HandlerFunc is an adapter to allow ordinary functions as handlers
type HandlerFunc func(ctx context.Context, j *job.Job) error

// Handle implements the Handler interface
func (f HandlerFunc) Handle(ctx context.Context, j *job.Job) error {
	return f(ctx, j)
}

// Middleware wraps a handler with additional functionality
type Middleware func(Handler) Handler

// Chain applies middlewares to a handler in order
func Chain(h Handler, middlewares ...Middleware) Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
