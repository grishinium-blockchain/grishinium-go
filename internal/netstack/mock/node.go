package mock

import (
	"context"
	"errors"
	"sync"

	"github.com/grishinium-blockchain/grishinium-go/internal/netstack"
)

// Node is a simple in-memory implementation of netstack.Node for bootstrap/testing.
type Node struct {
	cfg   netstack.Config
	mu    sync.RWMutex
	alive bool
	addr  string
	subs  map[string]chan []byte
	// dht-like state
	providers map[string][]string   // key -> list of provider addresses
	values    map[string][]byte     // key -> value
}

// Provide announces this node as a provider for the given key.
func (n *Node) Provide(ctx context.Context, key []byte) error {
    n.mu.Lock()
    defer n.mu.Unlock()
    k := string(key)
    list := n.providers[k]
    // avoid duplicates
    for _, a := range list {
        if a == n.addr {
            return nil
        }
    }
    n.providers[k] = append(list, n.addr)
    return nil
}

// FindProviders returns up to limit providers for the given key.
func (n *Node) FindProviders(ctx context.Context, key []byte, limit int) ([]string, error) {
    n.mu.RLock()
    defer n.mu.RUnlock()
    list := n.providers[string(key)]
    if limit <= 0 || limit > len(list) {
        limit = len(list)
    }
    out := make([]string, 0, limit)
    for i := 0; i < limit; i++ {
        out = append(out, list[i])
    }
    return out, nil
}

// PutValue stores a small value for the given key.
func (n *Node) PutValue(ctx context.Context, key, value []byte) error {
    n.mu.Lock()
    defer n.mu.Unlock()
    n.values[string(key)] = append([]byte(nil), value...)
    return nil
}

// GetValue retrieves a previously stored value for the key.
func (n *Node) GetValue(ctx context.Context, key []byte) ([]byte, error) {
    n.mu.RLock()
    defer n.mu.RUnlock()
    v, ok := n.values[string(key)]
    if !ok {
        return nil, nil
    }
    out := make([]byte, len(v))
    copy(out, v)
    return out, nil
}

func New(cfg netstack.Config) *Node {
	return &Node{cfg: cfg, addr: "mock://local", subs: make(map[string]chan []byte), providers: make(map[string][]string), values: make(map[string][]byte)}
}

func (n *Node) Start(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if n.alive {
		return nil
	}
	n.alive = true
	return nil
}

func (n *Node) Close(ctx context.Context) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.alive {
		return nil
	}
	for topic, ch := range n.subs {
		close(ch)
		delete(n.subs, topic)
	}
	n.alive = false
	return nil
}

// Addr returns a mock address string.
func (n *Node) Addr() string { return n.addr }

// PeerID returns a fixed mock peer ID.
func (n *Node) PeerID() string { return "mock-peer" }


func (n *Node) Publish(ctx context.Context, topic string, data []byte) error {
	n.mu.RLock()
	ch, ok := n.subs[topic]
	n.mu.RUnlock()
	if !ok {
		return errors.New("no subscribers")
	}
	select {
	case ch <- append([]byte(nil), data...):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (n *Node) Subscribe(ctx context.Context, topic string) (<-chan []byte, error) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if _, ok := n.subs[topic]; ok {
		return nil, errors.New("already subscribed")
	}
	ch := make(chan []byte, 1024)
	n.subs[topic] = ch
	out := make(chan []byte)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				select {
				case out <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out, nil
}

func (n *Node) Unsubscribe(ctx context.Context, topic string) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	ch, ok := n.subs[topic]
	if !ok {
		return nil
	}
	close(ch)
	delete(n.subs, topic)
	return nil
}

func (n *Node) FindPeer(ctx context.Context, id string) (string, error) {
	return "mock://peer/" + id, nil
}
