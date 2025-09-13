# GRISHINIUM Go Migration Scaffold

This directory mirrors the C++ project structure with Go packages and commands.
Each package contains placeholders to be implemented during migration.

Repository: https://github.com/grishinium-blockchain/grishinium-go

Build all commands under cmd/* with:
  go build ./...

Quick start (CLI):

```bash
go build -o bin/grishiniumlib-cli ./cmd/grishiniumlib-cli
./bin/grishiniumlib-cli -version
./bin/grishiniumlib-cli -endpoint 127.0.0.1:1234 -ping
```

Build tags

- By default, a lightweight in-memory mock networking stack is used (no extra deps).
- To enable the production libp2p-based networking stack, build with the tag `libp2p`:

```bash
go build -tags libp2p -o bin/validator-engine ./cmd/validator-engine
```

Storage

- In-memory KV is enabled by default (for dev/testing).
- A Pebble-based KV will be added under the build tag `pebble` (TBD).

Dependencies (to be fetched when ready)

- Networking: libp2p core, kad-dht, pubsub, multiaddr
- Crypto: secp256k1 (dcrd), bls12-381 (kilic)
- Storage: Pebble or Badger (TBD)

Once approved, we will fetch and pin these dependencies in go.mod and wire the implementations.
