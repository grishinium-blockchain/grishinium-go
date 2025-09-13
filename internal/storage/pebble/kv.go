//go:build pebble

package pebble

import (
	"context"
	"path/filepath"

	"github.com/cockroachdb/pebble"
	"github.com/grishinium-blockchain/grishinium-go/internal/storage"
)

// KV is a Pebble-backed implementation of storage.KV.
type KV struct {
	cfg storage.Config
	db  *pebble.DB
}

func New(cfg storage.Config) *KV { return &KV{cfg: cfg} }

func (kv *KV) Open(ctx context.Context) error {
	path := filepath.Clean(kv.cfg.Path)
	db, err := pebble.Open(path, &pebble.Options{})
	if err != nil {
		return err
	}
	kv.db = db
	return nil
}

func (kv *KV) Close(ctx context.Context) error {
	if kv.db == nil {
		return nil
	}
	return kv.db.Close()
}

func (kv *KV) Get(ctx context.Context, key []byte) ([]byte, error) {
	if kv.db == nil {
		return nil, nil
	}
	v, closer, err := kv.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	defer closer.Close()
	out := make([]byte, len(v))
	copy(out, v)
	return out, nil
}

func (kv *KV) Put(ctx context.Context, key, value []byte) error {
	if kv.db == nil {
		return nil
	}
	return kv.db.Set(key, value, pebble.NoSync)
}

func (kv *KV) Delete(ctx context.Context, key []byte) error {
	if kv.db == nil {
		return nil
	}
	return kv.db.Delete(key, pebble.NoSync)
}
