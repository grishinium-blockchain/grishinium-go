//go:build libp2p

package libp2p

import (
	"context"
	"fmt"
	"encoding/hex"

	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p/core/crypto"
	host "github.com/libp2p/go-libp2p/core/host"
	peer "github.com/libp2p/go-libp2p/core/peer"
	kad "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	ma "github.com/multiformats/go-multiaddr"
	cid "github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"

	"github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// Node is a libp2p-backed implementation of netstack.Node.
type Node struct {
	cfg    netstack.Config
	Host   host.Host
	DHT    *kad.IpfsDHT
	PubSub *pubsub.PubSub
}

// PeerID returns the string representation of the local host ID.
func (n *Node) PeerID() string {
    if n.Host == nil {
        return ""
    }
    return n.Host.ID().String()
}

func New(cfg netstack.Config) *Node { return &Node{cfg: cfg} }

func (n *Node) Start(ctx context.Context) error {
	// Build listen addrs
	var addrs []ma.Multiaddr
	for _, s := range n.cfg.ListenAddrs {
		m, err := ma.NewMultiaddr(s)
		if err != nil {
			return fmt.Errorf("invalid listen addr %q: %w", s, err)
		}
		addrs = append(addrs, m)
	}

	// Identity option (if provided)
	var opts []libp2p.Option
	opts = append(opts, libp2p.ListenAddrs(addrs...))
	if len(n.cfg.IdentityPriv) > 0 {
		priv, err := crypto.UnmarshalEd25519PrivateKey(n.cfg.IdentityPriv)
		if err != nil {
			return fmt.Errorf("unmarshal ed25519 identity: %w", err)
		}
		opts = append(opts, libp2p.Identity(priv))
	}

	h, err := libp2p.New(opts...)
	if err != nil {
		return err
	}
	n.Host = h

	// Create DHT
	dht, err := kad.New(ctx, h, kad.Mode(kad.ModeAuto))
	if err != nil {
		return err
	}
	n.DHT = dht
	if err := dht.Bootstrap(ctx); err != nil {
		return err
	}
	for _, b := range n.cfg.Bootstrap {
		addr, err := ma.NewMultiaddr(b)
		if err == nil {
			_ = n.connectAddr(ctx, addr)
		}
	}

	// PubSub
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		return err
	}
	n.PubSub = ps
	return nil
}

func (n *Node) connectAddr(ctx context.Context, addr ma.Multiaddr) error {
	pi, err := peerInfoFromAddr(addr)
	if err != nil {
		return err
	}
	return n.Host.Connect(ctx, *pi)
}

func (n *Node) Close(ctx context.Context) error {
	if n.DHT != nil {
		_ = n.DHT.Close()
	}
	if n.Host != nil {
		return n.Host.Close()
	}
	return nil
}

func (n *Node) Addr() string {
	if n.Host == nil {
		return ""
	}
	addrs := n.Host.Addrs()
	if len(addrs) == 0 {
		return ""
	}
	return addrs[0].String()
}

func (n *Node) Publish(ctx context.Context, topic string, data []byte) error {
	if n.PubSub == nil {
		return fmt.Errorf("pubsub not initialized")
	}
	t, err := n.PubSub.Join(topic)
	if err != nil {
		return err
	}
	defer t.Close()
	return t.Publish(ctx, data)
}

func (n *Node) Subscribe(ctx context.Context, topic string) (<-chan []byte, error) {
	if n.PubSub == nil {
		return nil, fmt.Errorf("pubsub not initialized")
	}
	t, err := n.PubSub.Join(topic)
	if err != nil {
		return nil, err
	}
	sub, err := t.Subscribe()
	if err != nil {
		_ = t.Close()
		return nil, err
	}
	out := make(chan []byte)
	go func() {
		defer close(out)
		defer sub.Cancel()
		defer t.Close()
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				return
			}
			select {
			case out <- append([]byte(nil), msg.Data...):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, nil
}

func (n *Node) Unsubscribe(ctx context.Context, topic string) error {
	// No-op: Subscribe manages its own lifecycle per call.
	return nil
}

func (n *Node) FindPeer(ctx context.Context, id string) (string, error) {
    if n.DHT == nil {
        return "", fmt.Errorf("dht not initialized")
    }
    // Placeholder: a real implementation converts id to peer.ID and queries DHT.
    return "libp2p://peer/" + id, nil
}

// peerInfoFromAddr parses a multiaddr with /p2p/peerID into a PeerInfo.
// NOTE: kept minimal to avoid pulling extra utils; implement in full later.
func peerInfoFromAddr(m ma.Multiaddr) (*peer.AddrInfo, error) {
    pi, err := peer.AddrInfoFromP2pAddr(m)
    if err != nil {
        return nil, err
    }
    return pi, nil
}

// --- DHT helpers and methods ---

func keyToCid(key []byte) (cid.Cid, error) {
    sum, err := mh.Sum(key, mh.SHA2_256, -1)
    if err != nil {
        return cid.Cid{}, err
    }
    return cid.NewCidV1(cid.Raw, sum), nil
}

func keyToDHTPath(key []byte) string {
    return "/grishinium/" + hex.EncodeToString(key)
}

// Provide announces this node as a provider for the given key via provider records.
func (n *Node) Provide(ctx context.Context, key []byte) error {
    if n.DHT == nil {
        return fmt.Errorf("dht not initialized")
    }
    c, err := keyToCid(key)
    if err != nil {
        return err
    }
    return n.DHT.Provide(ctx, c, true)
}

// FindProviders returns up to limit provider addresses for the given key.
func (n *Node) FindProviders(ctx context.Context, key []byte, limit int) ([]string, error) {
    if n.DHT == nil {
        return nil, fmt.Errorf("dht not initialized")
    }
    c, err := keyToCid(key)
    if err != nil {
        return nil, err
    }
    out := make([]string, 0, limit)
    ch := n.DHT.FindProvidersAsync(ctx, c, limit)
    for info := range ch {
        for _, a := range info.Addrs {
            out = append(out, a.String())
        }
        if limit > 0 && len(out) >= limit {
            break
        }
    }
    if limit > 0 && len(out) > limit {
        out = out[:limit]
    }
    return out, nil
}

// PutValue stores small metadata value under a namespaced key in DHT.
func (n *Node) PutValue(ctx context.Context, key, value []byte) error {
    if n.DHT == nil {
        return fmt.Errorf("dht not initialized")
    }
    k := keyToDHTPath(key)
    return n.DHT.PutValue(ctx, k, value)
}

// GetValue retrieves value stored under the namespaced key from DHT.
func (n *Node) GetValue(ctx context.Context, key []byte) ([]byte, error) {
    if n.DHT == nil {
        return nil, fmt.Errorf("dht not initialized")
    }
    k := keyToDHTPath(key)
    return n.DHT.GetValue(ctx, k)
}
