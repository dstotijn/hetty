// hetty-desktop launches Hetty in a native desktop window via Wails v3.
package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	llog "log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"go.uber.org/zap"

	"github.com/dstotijn/hetty/pkg/log"
	"github.com/dstotijn/hetty/pkg/server"
)

var version = "0.0.0"

//go:embed assets/icon.png
var iconPNG []byte

//go:embed admin
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminContent embed.FS

func main() {
	var (
		certPath = flag.String("cert", "~/.hetty/hetty_cert.pem", "Path to root CA certificate.")
		keyPath  = flag.String("key", "~/.hetty/hetty_key.pem", "Path to root CA private key.")
		dbPath   = flag.String("db", "~/.hetty/hetty.db", "Database file path.")
		addr     = flag.String("addr", "127.0.0.1:8080", "TCP address to listen on.")
		verbose  = flag.Bool("verbose", false, "Enable verbose logging.")
		jsonLogs = flag.Bool("json", false, "Encode logs as JSON.")
		showVer  = flag.Bool("version", false, "Output version.")
	)
	flag.BoolVar(showVer, "v", false, "Output version.")
	flag.Parse()

	if *showVer {
		fmt.Println(version)
		return
	}

	logger, err := log.NewZapLogger(*verbose, *jsonLogs)
	if err != nil {
		llog.Fatal(err)
	}
	defer logger.Sync() //nolint:errcheck

	fsSub, err := fs.Sub(adminContent, "admin")
	if err != nil {
		logger.Fatal("Failed to load embedded admin UI.", zap.Error(err))
	}

	inst, err := server.Start(server.Config{
		CertPath: *certPath,
		KeyPath:  *keyPath,
		DBPath:   *dbPath,
		Addr:     *addr,
		Logger:   logger,
		Version:  version,
		AdminFS:  fsSub,
	})
	if err != nil {
		logger.Fatal("Failed to start Hetty server.", zap.Error(err))
	}

	if err := server.WaitReady(inst.URL, 15*time.Second); err != nil {
		inst.Shutdown(context.Background()) //nolint:errcheck
		logger.Fatal("Hetty server did not become ready.", zap.Error(err))
	}

	mainApp := application.New(application.Options{
		Name:        "Hetty",
		Description: "HTTP toolkit for security research",
		Icon:        iconPNG,
		OnShutdown: func() {
			if err := inst.Shutdown(context.Background()); err != nil {
				logger.Error("Failed to shutdown Hetty server.", zap.Error(err))
			}
		},
	})

	mainApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:     "Hetty",
		Width:     1280,
		Height:    800,
		MinWidth:  1024,
		MinHeight: 600,
		URL:       inst.URL,
	})

	if err := mainApp.Run(); err != nil {
		logger.Fatal("Desktop application failed.", zap.Error(err))
	}
}
