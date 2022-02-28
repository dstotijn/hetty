package main

import (
	"context"
	"crypto/tls"
	"embed"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	llog "log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/chromedp/chromedp"
	badgerdb "github.com/dgraph-io/badger/v3"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dstotijn/hetty/pkg/api"
	"github.com/dstotijn/hetty/pkg/chrome"
	"github.com/dstotijn/hetty/pkg/db/badger"
	"github.com/dstotijn/hetty/pkg/log"
	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/sender"
)

var version = "0.0.0"

// Flag variables.
var (
	caCertFile   string
	caKeyFile    string
	dbPath       string
	addr         string
	launchChrome bool
	debug        bool
	noPrettyLogs bool
)

//go:embed admin
//go:embed admin/_next/static
//go:embed admin/_next/static/chunks/pages/*.js
//go:embed admin/_next/static/*/*.js
var adminContent embed.FS

func main() {
	ctx := context.Background()

	flag.StringVar(&caCertFile, "cert", "~/.hetty/hetty_cert.pem",
		"CA certificate filepath. Creates a new CA certificate if file doesn't exist")
	flag.StringVar(&caKeyFile, "key", "~/.hetty/hetty_key.pem",
		"CA private key filepath. Creates a new CA private key if file doesn't exist")
	flag.StringVar(&dbPath, "db", "~/.hetty/db", "Database directory path")
	flag.StringVar(&addr, "addr", ":8080", "TCP address to listen on, in the form \"host:port\"")
	flag.BoolVar(&launchChrome, "chrome", false, "Launch Chrome with proxy settings and certificate errors ignored")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&noPrettyLogs, "disable-pretty-logs", false, "Disable human readable console logs and encode with JSON.")
	flag.Parse()

	logger, err := log.NewZapLogger(debug, !noPrettyLogs)
	if err != nil {
		llog.Fatal(err)
	}
	defer logger.Sync()

	mainLogger := logger.Named("main")

	listenHost, listenPort, err := net.SplitHostPort(addr)
	if err != nil {
		mainLogger.Fatal("Failed to parse listening address.", zap.Error(err))
	}

	url := fmt.Sprintf("http://%v:%v", listenHost, listenPort)
	if listenHost == "" || listenHost == "0.0.0.0" || listenHost == "127.0.0.1" || listenHost == "::1" {
		url = fmt.Sprintf("http://localhost:%v", listenPort)
	}

	// Expand `~` in filepaths.
	caCertFile, err := homedir.Expand(caCertFile)
	if err != nil {
		logger.Fatal("Failed to parse CA certificate filepath.", zap.Error(err))
	}

	caKeyFile, err := homedir.Expand(caKeyFile)
	if err != nil {
		logger.Fatal("Failed to parse CA private key filepath.", zap.Error(err))
	}

	dbPath, err := homedir.Expand(dbPath)
	if err != nil {
		logger.Fatal("Failed to parse database path.", zap.Error(err))
	}

	// Load existing CA certificate and key from disk, or generate and write
	// to disk if no files exist yet.
	caCert, caKey, err := proxy.LoadOrCreateCA(caKeyFile, caCertFile)
	if err != nil {
		logger.Fatal("Failed to load or create CA key pair.", zap.Error(err))
	}

	// BadgerDB logs some verbose entries with `INFO` level, so unless
	// we're running in debug mode, bump the minimal level to `WARN`.
	dbLogger := logger.Named("badgerdb").WithOptions(zap.IncreaseLevel(zapcore.WarnLevel))

	dbSugaredLogger := dbLogger.Sugar()

	badger, err := badger.OpenDatabase(
		badgerdb.DefaultOptions(dbPath).WithLogger(badger.NewLogger(dbSugaredLogger)),
	)
	if err != nil {
		logger.Fatal("Failed to open database.", zap.Error(err))
	}
	defer badger.Close()

	scope := &scope.Scope{}

	reqLogService := reqlog.NewService(reqlog.Config{
		Scope:      scope,
		Repository: badger,
		Logger:     logger.Named("reqlog").Sugar(),
	})

	senderService := sender.NewService(sender.Config{
		Repository:    badger,
		ReqLogService: reqLogService,
	})

	projService, err := proj.NewService(proj.Config{
		Repository:    badger,
		ReqLogService: reqLogService,
		SenderService: senderService,
		Scope:         scope,
	})
	if err != nil {
		logger.Fatal("Failed to create new projects service.", zap.Error(err))
	}

	proxy, err := proxy.NewProxy(proxy.Config{
		CACert: caCert,
		CAKey:  caKey,
		Logger: logger.Named("proxy").Sugar(),
	})
	if err != nil {
		logger.Fatal("Failed to create new proxy.", zap.Error(err))
	}

	proxy.UseRequestModifier(reqLogService.RequestModifier)
	proxy.UseResponseModifier(reqLogService.ResponseModifier)

	fsSub, err := fs.Sub(adminContent, "admin")
	if err != nil {
		logger.Fatal("Failed to construct file system subtree from admin dir.", zap.Error(err))
	}

	adminHandler := http.FileServer(http.FS(fsSub))
	router := mux.NewRouter().SkipClean(true)
	adminRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		hostname, _ := os.Hostname()
		host, _, _ := net.SplitHostPort(req.Host)

		// Serve local admin routes when either:
		// - The `Host` is well-known, e.g. `hetty.proxy`, `localhost:[port]`
		//   or the listen addr `[host]:[port]`.
		// - The request is not for TLS proxying (e.g. no `CONNECT`) and not
		//   for proxying an external URL. E.g. Request-Line (RFC 7230, Section 3.1.1)
		//   has no scheme.
		return strings.EqualFold(host, hostname) ||
			req.Host == "hetty.proxy" ||
			req.Host == fmt.Sprintf("%v:%v", "localhost", listenPort) ||
			req.Host == fmt.Sprintf("%v:%v", listenHost, listenPort) ||
			req.Method != http.MethodConnect && !strings.HasPrefix(req.RequestURI, "http://")
	}).Subrouter().StrictSlash(true)

	// GraphQL server.
	gqlEndpoint := "/api/graphql/"
	adminRouter.Path(gqlEndpoint).Handler(api.HTTPHandler(&api.Resolver{
		ProjectService:    projService,
		RequestLogService: reqLogService,
		SenderService:     senderService,
	}, gqlEndpoint))

	// Admin interface.
	adminRouter.PathPrefix("").Handler(adminHandler)

	// Fallback (default) is the Proxy handler.
	router.PathPrefix("").Handler(proxy)

	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
		ErrorLog:     zap.NewStdLog(logger.Named("http")),
	}

	mainLogger.Info(fmt.Sprintf("Hetty (v%v) is running on %v ...", version, addr))
	mainLogger.Info(fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(32), "Get started at "+url))

	if launchChrome {
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

	err = httpServer.ListenAndServe()
	if err != nil && errors.Is(err, http.ErrServerClosed) {
		mainLogger.Fatal("HTTP server closed unexpected.", zap.Error(err))
	}
}
