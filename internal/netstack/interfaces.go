package netstack

import (
	"context"
)

// Node is a high-level GRISHINIUM networking node abstraction.
// It is intended to be backed by a professional-grade stack (libp2p: transport, Kad-DHT, PubSub).
// This interface allows swapping implementations for tests.
type Node interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error

	// Addr returns a human-readable address of the node (multiaddr or host:port).
	Addr() string

	// Publish broadcasts data to a topic within the overlay network.
	Publish(ctx context.Context, topic string, data []byte) error
	// Subscribe subscribes to a topic and returns a receive-only channel with messages.
	Subscribe(ctx context.Context, topic string) (<-chan []byte, error)
	// Unsubscribe closes a previously created subscription.
	Unsubscribe(ctx context.Context, topic string) error

	// FindPeer locates a peer by ID using DHT.
	FindPeer(ctx context.Context, id string) (string, error)
}

// Config configures the networking node.
type Config struct {
	ListenAddrs []string // e.g., "/ip4/0.0.0.0/tcp/0"
	Bootstrap   []string // bootstrap peers (multiaddrs)
}
