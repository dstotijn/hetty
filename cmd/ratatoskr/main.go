// Command ratatoskr is the Ratatoskr web security proxy.
//
// It currently serves a placeholder page; the proxy, the store and the Svelte UI
// are wired in later tasks (docs/08-roadmap.md).
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

const (
	defaultListen = "127.0.0.1:8080"

	readHeaderTimeout = 10 * time.Second
	idleTimeout       = 60 * time.Second
	maxHeaderBytes    = 1 << 20
	shutdownGrace     = 5 * time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, os.Args[1:]); err != nil {
		slog.Error("ratatoskr stopped", "error", err)
		os.Exit(1)
	}
}

// options is the CLI surface of cmd/ratatoskr. internal/config replaces it in
// T-004, which is also where the data dir gets used.
type options struct {
	listen      string
	dataDir     string
	showVersion bool
}

func parseFlags(args []string) (options, error) {
	var opts options

	fs := flag.NewFlagSet("ratatoskr", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // flag errors are returned and logged once, by main.
	fs.StringVar(&opts.listen, "listen", defaultListen, "address to serve the UI/API on")
	fs.StringVar(&opts.dataDir, "data-dir", "", "directory for config, CA and project databases")
	fs.BoolVar(&opts.showVersion, "version", false, "print the version and exit")

	if err := fs.Parse(args); err != nil {
		return options{}, fmt.Errorf("parse flags: %w", err)
	}

	return opts, nil
}

func run(ctx context.Context, args []string) error {
	opts, err := parseFlags(args)
	if err != nil {
		return err
	}

	if opts.showVersion {
		if _, err := fmt.Fprintf(os.Stdout, "ratatoskr %s\n", version); err != nil {
			return fmt.Errorf("write version: %w", err)
		}

		return nil
	}

	ln, err := net.Listen("tcp", opts.listen)
	if err != nil {
		return fmt.Errorf("listen on %q: %w", opts.listen, err)
	}

	return serve(ctx, ln, opts)
}

// serve runs the HTTP server on ln until ctx is cancelled, then drains in-flight
// requests. The listener is opened by the caller so the bound address is known
// (and logged) even for port 0.
func serve(ctx context.Context, ln net.Listener, opts options) error {
	attrs := []any{"addr", ln.Addr().String(), "version", version}
	if opts.dataDir != "" {
		attrs = append(attrs, "data_dir", opts.dataDir)
	}

	slog.Info("ratatoskr listening", attrs...)

	srv := &http.Server{
		Handler:           newHandler(),
		ReadHeaderTimeout: readHeaderTimeout,
		IdleTimeout:       idleTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	select {
	case err := <-serveErr:
		return fmt.Errorf("serve: %w", err)
	case <-ctx.Done():
	}

	// The signal context is already cancelled; give the drain its own deadline.
	shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownGrace)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	slog.Info("ratatoskr stopped")

	return nil
}
