package dht

import "context"

// Key is a DHT key (e.g., hash or content address).
type Key []byte

// Value is a DHT stored value payload.
type Value []byte

// Peer describes a DHT peer.
type Peer struct {
	ID   string
	Addr string // multiaddr or host:port, implementation-defined
}

// Table is the high-level DHT interface.
type Table interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
	Self() Peer

	// FindPeer locates a peer by ID.
	FindPeer(ctx context.Context, id string) (Peer, error)
	// FindProviders returns peers providing the given key/value.
	FindProviders(ctx context.Context, key Key, limit int) ([]Peer, error)
	// Provide announces that this node can provide the given key/value.
	Provide(ctx context.Context, key Key) error
	// Get/Put for small metadata values.
	Get(ctx context.Context, key Key) (Value, error)
	Put(ctx context.Context, key Key, value Value) error
}
