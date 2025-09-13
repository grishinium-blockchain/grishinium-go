package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "time"

    appctx "github.com/grishinium-blockchain/grishinium-go/internal/appctx"
    cfgpkg "github.com/grishinium-blockchain/grishinium-go/internal/config"
    logger "github.com/grishinium-blockchain/grishinium-go/internal/log"
)

func usage() {
    fmt.Fprintf(os.Stderr, "validator-engine\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  validator-engine -endpoint <host:port> [-debug] [-timeout 10s]\n\n")
    flag.PrintDefaults()
}

func main() {
    var cfg cfgpkg.Config

    cfgpkg.Flags(nil, &cfg)
    flag.Usage = usage
    flag.Parse()

    if cfg.Debug {
        logger.SetDebug()
    }

    root, cancel := appctx.WithSignals(context.Background())
    defer cancel()

    if cfg.Timeout <= 0 {
        cfg.Timeout = 10 * time.Second
    }
    ctx, opCancel := context.WithTimeout(root, cfg.Timeout)
    defer opCancel()

    // TODO: initialize production-grade networking (libp2p: transport, kad-dht, pubsub)
    // TODO: initialize storage (Pebble/Badger) and state
    // TODO: initialize validator services and start event loop

    // Temporary output while skeleton is in place
    fmt.Println("GRISHINIUM validator-engine skeleton: OK (services not yet implemented)")

    // Block until context cancellation (signals)
    <-ctx.Done()
}
