package main

import (
	"github.com/grishinium-blockchain/grishinium-go/internal/netstack"
	mocknet "github.com/grishinium-blockchain/grishinium-go/internal/netstack/mock"
)

// newNetstackNode returns the default (mock) implementation when built without tags.
func newNetstackNode(cfg netstack.Config) netstack.Node {
	return mocknet.New(cfg)
}
