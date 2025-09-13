package dht

import (
	"context"
	"errors"

	ns "github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// Adapter bridges dht.Table к реализации поверх internal/netstack.Node.
// Реализует только доступные операции; остальное помечено TODO.
type Adapter struct {
	node ns.Node
	self Peer
}

func NewAdapter(node ns.Node, self Peer) *Adapter { return &Adapter{node: node, self: self} }

func (a *Adapter) Start(ctx context.Context) error { return nil }
func (a *Adapter) Close(ctx context.Context) error { return nil }
func (a *Adapter) Self() Peer                        { return a.self }

func (a *Adapter) FindPeer(ctx context.Context, id string) (Peer, error) {
	addr, err := a.node.FindPeer(ctx, id)
	if err != nil {
		return Peer{}, err
	}
	return Peer{ID: id, Addr: addr}, nil
}

func (a *Adapter) FindProviders(ctx context.Context, key Key, limit int) ([]Peer, error) {
	addrs, err := a.node.FindProviders(ctx, key, limit)
	if err != nil {
		return nil, err
	}
	peers := make([]Peer, 0, len(addrs))
	for _, addr := range addrs {
		peers = append(peers, Peer{Addr: addr})
	}
	return peers, nil
}

func (a *Adapter) Provide(ctx context.Context, key Key) error {
	return a.node.Provide(ctx, key)
}

func (a *Adapter) Get(ctx context.Context, key Key) (Value, error) {
	return a.node.GetValue(ctx, key)
}

func (a *Adapter) Put(ctx context.Context, key Key, value Value) error {
	return a.node.PutValue(ctx, key, value)
}
