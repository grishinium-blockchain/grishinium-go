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
}

func New(cfg netstack.Config) *Node {
	return &Node{cfg: cfg, addr: "mock://local", subs: make(map[string]chan []byte)}
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

func (n *Node) Addr() string { return n.addr }

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
