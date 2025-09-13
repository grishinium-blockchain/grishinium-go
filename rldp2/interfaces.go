package rldp2

import "context"

// Frame represents a framed piece of data in RLDPv2.
type Frame []byte

// Stream provides a duplex byte stream over RLDPv2.
type Stream interface {
	ID() string
	Write(ctx context.Context, p []byte) (int, error)
	Read(ctx context.Context, buf []byte) (int, error)
	Close(ctx context.Context) error
}

// Manager manages RLDPv2 streams and sessions.
type Manager interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Open(ctx context.Context, id string) (Stream, error)
	Accept(ctx context.Context) (Stream, error)
}
