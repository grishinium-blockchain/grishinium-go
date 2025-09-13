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
	// TODO: интеграция с Kad-DHT: Providers.
	return nil, errors.New("FindProviders: not implemented yet")
}

func (a *Adapter) Provide(ctx context.Context, key Key) error {
	// TODO: интеграция с Kad-DHT: Provide.
	return errors.New("Provide: not implemented yet")
}

func (a *Adapter) Get(ctx context.Context, key Key) (Value, error) {
	// TODO: интеграция с Kad-DHT: Get value (record).
	return nil, errors.New("Get: not implemented yet")
}

func (a *Adapter) Put(ctx context.Context, key Key, value Value) error {
	// TODO: интеграция с Kad-DHT: Put value (record).
	return errors.New("Put: not implemented yet")
}
