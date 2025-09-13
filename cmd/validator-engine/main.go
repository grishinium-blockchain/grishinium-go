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
    overlaypkg "github.com/grishinium-blockchain/grishinium-go/overlay"
    dhtpkg "github.com/grishinium-blockchain/grishinium-go/dht"
)

func usage() {
    fmt.Fprintf(os.Stderr, "validator-engine\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  validator-engine -endpoint <host:port> [-listen <multiaddr>] [-bootstrap <multiaddr>] [-identity <path>] [-mdns] [-debug] [-timeout 10s]\n\n")
    flag.PrintDefaults()
}

func main() {
    var cfg cfgpkg.Config
    var listen multiFlag
    var bootstrap multiFlag
    var identityPath string
    var enableMDNS bool

    cfgpkg.Flags(nil, &cfg)
    flag.Var(&listen, "listen", "Listen multiaddr (repeatable). Example: /ip4/0.0.0.0/tcp/0")
    flag.Var(&bootstrap, "bootstrap", "Bootstrap peer multiaddr with /p2p/<peerID> (repeatable)")
    flag.StringVar(&identityPath, "identity", "", "Path to ed25519 private key file (raw). If empty, ephemeral identity is used")
    flag.BoolVar(&enableMDNS, "mdns", false, "Enable LAN peer discovery via mDNS (libp2p builds)")
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
    if enableMDNS {
        if err := ns.EnableMDNS(context.Background()); err != nil {
            fmt.Fprintln(os.Stderr, "mdns enable warning:", err)
        }
    }
    if addr := ns.Addr(); addr != "" {
        fmt.Println("listening:", addr)
        if pid := ns.PeerID(); pid != "" {
            fmt.Println("bootstrap:", addr+"/p2p/"+pid)
        }
    }

    // Wire adapters for overlay and DHT
    ov := overlaypkg.NewAdapter(ns)
    table := dhtpkg.NewAdapter(ns, dhtpkg.Peer{ID: ns.PeerID(), Addr: ns.Addr()})
    _ = ov
    _ = table

    // Example: subscribe to a default broadcast topic in debug mode
    if cfg.Debug {
        topic := overlaypkg.Topic("grishinium.broadcast")
        ch, err := ov.Subscribe(root, topic)
        if err == nil {
            go func() {
                for {
                    select {
                    case msg, ok := <-ch:
                        if !ok {
                            return
                        }
                        fmt.Println("overlay msg:", string(msg))
                    case <-root.Done():
                        return
                    }
                }
            }()
            // Publish a hello message
            _ = ov.Publish(context.Background(), topic, []byte("hello from validator-engine"))
        }
        // DHT demo: Put/Get and Provide/FindProviders
        key := []byte("demo-key")
        val := []byte("demo-value")
        _ = table.Put(root, key, val)
        if got, err := table.Get(root, key); err == nil && len(got) > 0 {
            fmt.Println("dht get:", string(got))
        }
        _ = table.Provide(root, key)
        if peers, err := table.FindProviders(root, key, 5); err == nil {
            fmt.Printf("dht providers: %d\n", len(peers))
        }
    }

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
