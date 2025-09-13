package grishiniumlib

import (
	"context"
	"errors"
	"net"
	"time"
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
	if c.cfg.Endpoint == "" {
		return errors.New("endpoint is empty")
	}
	d := &net.Dialer{}
	// derive timeout from context if set, otherwise use a small default
	timeout := 3 * time.Second
	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline)
		if timeout <= 0 {
			timeout = 3 * time.Second
		}
	}
	conn, err := d.DialContext(context.WithTimeout(ctx, timeout), "tcp", c.cfg.Endpoint)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func (c *stubClient) GetVersion(ctx context.Context) (string, error) {
	return "grishiniumlib " + iv.Version, nil
}
