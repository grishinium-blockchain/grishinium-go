package rldp

import "context"

// Chunk is a data fragment in RLDP transfers.
type Chunk []byte

// Sender sends data over RLDP.
type Sender interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Send(ctx context.Context, streamID string, data []byte) error
}

// Receiver receives data over RLDP.
type Receiver interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Recv(ctx context.Context, streamID string) ([]byte, error)
}

// Session coordinates a bidirectional RLDP transfer session.
type Session interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	OpenStream(ctx context.Context, id string) error
	CloseStream(ctx context.Context, id string) error
}
