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
    netstack "github.com/grishinium-blockchain/grishinium-go/internal/netstack"
    "github.com/grishinium-blockchain/grishinium-go/keyring"
)

func usage() {
    fmt.Fprintf(os.Stderr, "validator-engine\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  validator-engine -endpoint <host:port> [-listen <multiaddr>] [-bootstrap <multiaddr>] [-debug] [-timeout 10s]\n\n")
    flag.PrintDefaults()
}

func main() {
    var cfg cfgpkg.Config
    var listen multiFlag
    var bootstrap multiFlag
    var identityPath string

    cfgpkg.Flags(nil, &cfg)
    flag.Var(&listen, "listen", "Listen multiaddr (repeatable). Example: /ip4/0.0.0.0/tcp/0")
    flag.Var(&bootstrap, "bootstrap", "Bootstrap peer multiaddr with /p2p/<peerID> (repeatable)")
    flag.StringVar(&identityPath, "identity", "", "Path to ed25519 private key file (raw). If empty, ephemeral identity is used")
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

    // Load or generate identity
    id, err := keyring.LoadIdentity(identityPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "identity load error:", err)
        os.Exit(1)
    }

    // Build netstack config
    nsCfg := netstack.Config{ListenAddrs: listen, Bootstrap: bootstrap, IdentityPriv: []byte(id.Private)}
    var ns netstack.Node = newNetstackNode(nsCfg)
    if err := ns.Start(ctx); err != nil {
        fmt.Fprintln(os.Stderr, "netstack start error:", err)
        os.Exit(1)
    }
    defer ns.Close(context.Background())

    // TODO: initialize storage (Pebble/Badger) and state
    // TODO: initialize validator services and start event loop

    // Temporary output while skeleton is in place
    fmt.Println("GRISHINIUM validator-engine skeleton: OK (services not yet implemented)")
    if fp := id.Fingerprint(); fp != "" {
        fmt.Println("identity:", fp)
    }

    // Block until context cancellation (signals)
    <-ctx.Done()
}
