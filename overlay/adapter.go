package overlay

import (
	"context"

	ns "github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// Adapter bridges overlay.Manager to a netstack.Node implementation.
// Lifecycle (Start/Close) is managed by the underlying node, so Adapter methods are no-ops.
type Adapter struct {
	node ns.Node
}

func NewAdapter(node ns.Node) *Adapter { return &Adapter{node: node} }

func (a *Adapter) Start(ctx context.Context) error { return nil }
func (a *Adapter) Close(ctx context.Context) error { return nil }

func (a *Adapter) Publish(ctx context.Context, topic Topic, data []byte) error {
	return a.node.Publish(ctx, string(topic), data)
}

func (a *Adapter) Subscribe(ctx context.Context, topic Topic) (<-chan []byte, error) {
	return a.node.Subscribe(ctx, string(topic))
}

func (a *Adapter) Unsubscribe(ctx context.Context, topic Topic) error {
	return a.node.Unsubscribe(ctx, string(topic))
}

func (a *Adapter) Join(ctx context.Context, overlayID string) error {
	// Mapping overlays to topics can be done via naming convention.
	// Here we treat overlayID as a topic namespace; no extra action required.
	return nil
}

func (a *Adapter) Leave(ctx context.Context, overlayID string) error {
	// Same as Join; no persistent membership in this adapter.
	return nil
}
