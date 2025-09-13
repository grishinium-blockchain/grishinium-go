package grishiniumlib

import (
	"context"
	iv "github.com/grishinium-blockchain/grishinium-go/internal/version"
)

// Client defines the high-level GRISHINIUM client interface.
// Methods are placeholders and will be implemented during migration.
type Client interface {
	// Ping checks connectivity to a GRISHINIUM endpoint.
	Ping(ctx context.Context) error

	// GetVersion returns the node/library version information.
	GetVersion(ctx context.Context) (string, error)
}

// Config holds minimal client configuration.
type Config struct {
	Endpoint string // e.g. host:port
	APIKey   string // optional
}

// NewClient creates a new GRISHINIUM client instance.
// Currently returns a stub implementation.
func NewClient(cfg Config) Client {
	return &stubClient{cfg: cfg}
}

type stubClient struct {
	cfg Config
}

func (c *stubClient) Ping(ctx context.Context) error {
	return nil
}

func (c *stubClient) GetVersion(ctx context.Context) (string, error) {
    return "grishiniumlib " + iv.Version, nil
}
