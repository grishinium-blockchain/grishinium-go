package main

import (
    "context"
    "flag"
    "fmt"
    "net/http"
    "os"
    "time"

    appctx "github.com/grishinium-blockchain/grishinium-go/internal/appctx"
    logger "github.com/grishinium-blockchain/grishinium-go/internal/log"
)

func usage() {
    fmt.Fprintf(os.Stderr, "rldp-http-proxy\n")
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "  rldp-http-proxy -listen :8080 -backend grishinium://host:port [-debug]\n\n")
    flag.PrintDefaults()
}

func main() {
    var (
        listen  string
        backend string
        debug   bool
    )

    flag.StringVar(&listen, "listen", ":8080", "HTTP listen address")
    flag.StringVar(&backend, "backend", "", "GRISHINIUM backend (rldp:// or grishinium:// endpoint)")
    flag.BoolVar(&debug, "debug", false, "enable debug logging")
    flag.Usage = usage
    flag.Parse()

    if debug {
        logger.SetDebug()
    }

    root, cancel := appctx.WithSignals(context.Background())
    defer cancel()

    srv := &http.Server{Addr: listen, Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: implement RLDP bridge to GRISHINIUM backend
        w.WriteHeader(http.StatusNotImplemented)
        _, _ = w.Write([]byte("rldp-http-proxy: TODO implement RLDP bridge"))
    })}

    go func() {
        _ = srv.ListenAndServe()
    }()

    <-root.Done()
    ctx, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel2()
    _ = srv.Shutdown(ctx)
}
