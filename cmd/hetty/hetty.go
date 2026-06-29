package main

import (
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"

	"github.com/chromedp/chromedp"
	"github.com/peterbourgon/ff/v3/ffcli"
	"go.uber.org/zap"

	"github.com/dstotijn/hetty/pkg/chrome"
	"github.com/dstotijn/hetty/pkg/server"
)

var version = "0.0.0"

//go:embed admin
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminContent embed.FS

var hettyUsage = `
Usage:
    hetty [flags] [subcommand] [flags]

Runs an HTTP server with (MITM) proxy, GraphQL service, and a web based admin interface.

Options:
    --cert         Path to root CA certificate. Creates file if it doesn't exist. (Default: "~/.hetty/hetty_cert.pem")
    --key          Path to root CA private key. Creates file if it doesn't exist. (Default: "~/.hetty/hetty_key.pem")
    --db           Database file path. Creates file if it doesn't exist. (Default: "~/.hetty/hetty.db")
    --addr         TCP address for HTTP server to listen on, in the form \"host:port\". (Default: ":8080")
    --chrome       Launch Chrome with proxy settings applied and certificate errors ignored. (Default: false)
    --verbose      Enable verbose logging.
    --json         Encode logs as JSON, instead of pretty/human readable output.
    --version, -v  Output version.
    --help, -h     Output this usage text.

Subcommands:
    - cert  Certificate management

Run ` + "`hetty <subcommand> --help`" + ` for subcommand specific usage instructions.

Visit https://hetty.xyz to learn more about Hetty.
`

type HettyCommand struct {
	config *Config

	cert    string
	key     string
	db      string
	addr    string
	chrome  bool
	version bool
}

func NewHettyCommand() (*ffcli.Command, *Config) {
	cmd := HettyCommand{
		config: &Config{},
	}

	fs := flag.NewFlagSet("hetty", flag.ExitOnError)

	fs.StringVar(&cmd.cert, "cert", "~/.hetty/hetty_cert.pem",
		"Path to root CA certificate. Creates a new certificate if file doesn't exist.")
	fs.StringVar(&cmd.key, "key", "~/.hetty/hetty_key.pem",
		"Path to root CA private key. Creates a new private key if file doesn't exist.")
	fs.StringVar(&cmd.db, "db", "~/.hetty/hetty.db", "Database file path. Creates file if it doesn't exist.")
	fs.StringVar(&cmd.addr, "addr", ":8080", "TCP address to listen on, in the form \"host:port\".")
	fs.BoolVar(&cmd.chrome, "chrome", false, "Launch Chrome with proxy settings applied and certificate errors ignored.")
	fs.BoolVar(&cmd.version, "version", false, "Output version.")
	fs.BoolVar(&cmd.version, "v", false, "Output version.")

	cmd.config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:    "hetty",
		FlagSet: fs,
		Subcommands: []*ffcli.Command{
			NewCertCommand(cmd.config),
		},
		Exec: cmd.Exec,
		UsageFunc: func(*ffcli.Command) string {
			return hettyUsage
		},
	}, cmd.config
}

func (cmd *HettyCommand) Exec(ctx context.Context, _ []string) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	if cmd.version {
		fmt.Fprint(os.Stdout, version+"\n")
		return nil
	}

	mainLogger := cmd.config.logger.Named("main")

	fsSub, err := fs.Sub(adminContent, "admin")
	if err != nil {
		cmd.config.logger.Fatal("Failed to construct file system subtree from admin dir.", zap.Error(err))
	}

	inst, err := server.Start(server.Config{
		CertPath: cmd.cert,
		KeyPath:  cmd.key,
		DBPath:   cmd.db,
		Addr:     cmd.addr,
		Logger:   cmd.config.logger,
		Version:  version,
		AdminFS:  fsSub,
	})
	if err != nil {
		cmd.config.logger.Fatal("Failed to start Hetty server.", zap.Error(err))
	}

	url := inst.URL

	if cmd.chrome {
		ctx, cancel := chrome.NewExecAllocator(ctx, chrome.Config{
			ProxyServer:      url,
			ProxyBypassHosts: []string{url},
		})
		defer cancel()

		taskCtx, cancel := chromedp.NewContext(ctx)
		defer cancel()

		err = chromedp.Run(taskCtx, chromedp.Navigate(url))

		switch {
		case errors.Is(err, exec.ErrNotFound):
			mainLogger.Info("Chrome executable not found.")
		case err != nil:
			mainLogger.Error(fmt.Sprintf("Failed to navigate to %v.", url), zap.Error(err))
		default:
			mainLogger.Info("Launched Chrome.")
		}
	}

	<-ctx.Done()
	stop()

	mainLogger.Info("Shutting down HTTP server. Press Ctrl+C to force quit.")

	err = inst.Shutdown(context.Background())
	if err != nil {
		return fmt.Errorf("failed to shutdown Hetty server: %w", err)
	}

	return nil
}
