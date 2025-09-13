package rldp2

import (
	"context"
	"errors"

	ns "github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// ManagerImpl is a skeleton RLDPv2 manager that will use libp2p streams under the hood.
type ManagerImpl struct {
	node ns.Node
	alive bool
}

// NewManager creates a new RLDPv2 manager bound to the provided networking node.
func NewManager(node ns.Node) *ManagerImpl { return &ManagerImpl{node: node} }

func (m *ManagerImpl) Start(ctx context.Context) error {
	if m.alive {
		return nil
	}
	m.alive = true
	return nil
}

func (m *ManagerImpl) Close(ctx context.Context) error {
	if !m.alive {
		return nil
	}
	m.alive = false
	return nil
}

func (m *ManagerImpl) Open(ctx context.Context, id string) (Stream, error) {
	return nil, errors.New("rldp2.ManagerImpl.Open: not implemented yet")
}

func (m *ManagerImpl) Accept(ctx context.Context) (Stream, error) {
	return nil, errors.New("rldp2.ManagerImpl.Accept: not implemented yet")
}
