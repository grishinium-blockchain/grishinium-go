package config

import (
	"flag"
	"time"
)

// Config is a base configuration used across binaries.
type Config struct {
	Endpoint string        // network endpoint in form host:port
	Timeout  time.Duration // default request timeout
	Debug    bool          // enable debug logging
}

// Flags binds common flags into the provided FlagSet (or flag.CommandLine when nil).
func Flags(fs *flag.FlagSet, cfg *Config) {
	if fs == nil {
		fs = flag.CommandLine
	}
	fs.StringVar(&cfg.Endpoint, "endpoint", "", "GRISHINIUM endpoint in form host:port")
	fs.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "default request timeout")
	fs.BoolVar(&cfg.Debug, "debug", false, "enable debug logging")
}
