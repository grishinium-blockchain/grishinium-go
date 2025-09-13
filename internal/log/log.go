package log

import (
	"log"
	"log/slog"
	"os"
)

var (
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
)

// SetDebug enables verbose logging.
func SetDebug() {
	Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.SetOutput(os.Stderr)
}
