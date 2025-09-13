package tdactor

import "context"

// Mailbox is a minimal mailbox abstraction for actor messages.
type Mailbox[T any] chan T

// Actor represents a basic actor with lifecycle hooks.
type Actor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

// Supervisor can supervise a group of actors.
type Supervisor interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
