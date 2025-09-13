package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "time"

    "github.com/grishinium-blockchain/grishinium-go/grishiniumlib"
)

func usage() {
    fmt.Fprintf(os.Stderr, "grishiniumlib-cli\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  grishiniumlib-cli -endpoint <host:port> [-ping] [-version]\n\n")
    flag.PrintDefaults()
}

func main() {
    var (
        endpoint string
        doPing   bool
        showVer  bool
        timeout  time.Duration
    )

    flag.StringVar(&endpoint, "endpoint", "", "GRISHINIUM endpoint in the form host:port")
    flag.BoolVar(&doPing, "ping", false, "Perform a connectivity ping")
    flag.BoolVar(&showVer, "version", false, "Print library version")
    flag.DurationVar(&timeout, "timeout", 5*time.Second, "Request timeout")
    flag.Usage = usage
    flag.Parse()

    if !doPing && !showVer {
        usage()
        os.Exit(2)
    }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    client := grishiniumlib.NewClient(grishiniumlib.Config{Endpoint: endpoint})

    if showVer {
        ver, err := client.GetVersion(ctx)
        if err != nil {
            fmt.Fprintln(os.Stderr, "error:", err)
            os.Exit(1)
        }
        fmt.Println(ver)
    }

    if doPing {
        if endpoint == "" {
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
