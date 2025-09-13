//go:build libp2p

package main

import (
	"github.com/grishinium-blockchain/grishinium-go/internal/netstack"
	libp2pnode "github.com/grishinium-blockchain/grishinium-go/internal/netstack/libp2p"
)

// newNetstackNode returns the libp2p implementation when built with -tags libp2p.
func newNetstackNode(cfg netstack.Config) netstack.Node {
	return libp2pnode.New(cfg)
}
