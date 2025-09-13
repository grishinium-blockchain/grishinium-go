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
    "github.com/grishinium-blockchain/grishinium-go/grishiniumlib"
)

func usage() {
    fmt.Fprintf(os.Stderr, "lite-client\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  lite-client -endpoint <host:port> [-ping] [-version] [-debug] [-timeout 5s]\n\n")
    flag.PrintDefaults()
}

func main() {
    var (
        cfg cfgpkg.Config
        doPing  bool
        showVer bool
    )

    cfgpkg.Flags(nil, &cfg)
    flag.BoolVar(&doPing, "ping", false, "Ping GRISHINIUM endpoint and exit")
    flag.BoolVar(&showVer, "version", false, "Print client/library version and exit")
    flag.Usage = usage
    flag.Parse()

    if cfg.Debug {
        logger.SetDebug()
    }

    if !doPing && !showVer {
        usage()
        os.Exit(2)
    }

    root, cancel := appctx.WithSignals(context.Background())
    defer cancel()

    if cfg.Timeout <= 0 {
        cfg.Timeout = 5 * time.Second
    }
    ctx, opCancel := context.WithTimeout(root, cfg.Timeout)
    defer opCancel()

    client := grishiniumlib.NewClient(grishiniumlib.Config{Endpoint: cfg.Endpoint})

    if showVer {
        ver, err := client.GetVersion(ctx)
        if err != nil {
            fmt.Fprintln(os.Stderr, "error:", err)
            os.Exit(1)
        }
        fmt.Println(ver)
        if !doPing {
            return
        }
    }

    if doPing {
        if cfg.Endpoint == "" {
            fmt.Fprintln(os.Stderr, "-endpoint is required for -ping")
            os.Exit(2)
        }
        if err := client.Ping(ctx); err != nil {
            fmt.Fprintln(os.Stderr, "ping failed:", err)
            os.Exit(1)
        }
        fmt.Println("ping: OK")
    }
}
