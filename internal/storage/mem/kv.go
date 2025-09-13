package mem

import (
	"context"
	"sync"

	"github.com/grishinium-blockchain/grishinium-go/internal/storage"
)

// KV is an in-memory implementation of storage.KV for development and tests.
type KV struct {
	cfg storage.Config
	mu  sync.RWMutex
	m   map[string][]byte
}

func New(cfg storage.Config) *KV { return &KV{cfg: cfg, m: make(map[string][]byte)} }

func (kv *KV) Open(ctx context.Context) error { return nil }
func (kv *KV) Close(ctx context.Context) error { return nil }

func (kv *KV) Get(ctx context.Context, key []byte) ([]byte, error) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	v, ok := kv.m[string(key)]
	if !ok {
		return nil, nil
	}
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (kv *KV) Put(ctx context.Context, key, value []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.m[string(key)] = append([]byte(nil), value...)
	return nil
}

func (kv *KV) Delete(ctx context.Context, key []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	delete(kv.m, string(key))
	return nil
}
