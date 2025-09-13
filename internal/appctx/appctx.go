package appctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// WithSignals returns a context canceled on SIGINT/SIGTERM and a cancel func.
func WithSignals(parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	ch := make(chan os.Signal, 2)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			cancel()
		}
	}()
	return ctx, cancel
}
