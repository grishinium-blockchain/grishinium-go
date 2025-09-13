package rldp

import (
	"context"
	"errors"

	ns "github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// ManagerImpl is a skeleton RLDP manager that will use libp2p streams under the hood.
type ManagerImpl struct {
	node ns.Node
	alive bool
}

// NewManager creates a new RLDP manager bound to the provided networking node.
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

func (m *ManagerImpl) Send(ctx context.Context, streamID string, data []byte) error {
	return errors.New("rldp.ManagerImpl.Send: not implemented yet")
}

// Ensure ManagerImpl satisfies the Sender/Receiver interfaces if needed.
var _ Sender = (*ManagerImpl)(nil)
var _ Receiver = (*ManagerImpl)(nil)
