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
