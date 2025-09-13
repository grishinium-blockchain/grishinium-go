package storage

import "context"

// KV is a minimal key-value storage interface suitable for validator state.
type KV interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error
	Get(ctx context.Context, key []byte) ([]byte, error)
	Put(ctx context.Context, key, value []byte) error
	Delete(ctx context.Context, key []byte) error
}

// Config describes storage configuration.
type Config struct {
	Path string
	ReadOnly bool
}
