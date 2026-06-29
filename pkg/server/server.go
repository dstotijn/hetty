package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"

	"github.com/dstotijn/hetty/pkg/api"
	"github.com/dstotijn/hetty/pkg/db/bolt"
	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/proxy/intercept"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/sender"
)

// Config configures the Hetty HTTP server (proxy, GraphQL API, and admin UI).
type Config struct {
	CertPath string
	KeyPath  string
	DBPath   string
	Addr     string
	Logger   *zap.Logger
	Version  string
	AdminFS  fs.FS
}

// Instance is a running Hetty server.
type Instance struct {
	HTTPServer *http.Server
	URL        string
	boltDB     *bolt.Database
	logger     *zap.Logger
}

// Start boots the Hetty HTTP server and returns before blocking on requests.
func Start(cfg Config) (*Instance, error) {
	mainLogger := cfg.Logger.Named("main")

	listenHost, listenPort, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("parse listen address: %w", err)
	}

	url := fmt.Sprintf("http://%v:%v", listenHost, listenPort)
	if listenHost == "" || listenHost == "0.0.0.0" || listenHost == "127.0.0.1" || listenHost == "::1" {
		url = fmt.Sprintf("http://localhost:%v", listenPort)
	}

	caCertFile, err := homedir.Expand(cfg.CertPath)
	if err != nil {
		return nil, fmt.Errorf("expand CA certificate path: %w", err)
	}

	caKeyFile, err := homedir.Expand(cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("expand CA private key path: %w", err)
	}

	dbPath, err := homedir.Expand(cfg.DBPath)
	if err != nil {
		return nil, fmt.Errorf("expand database path: %w", err)
	}

	caCert, caKey, err := proxy.LoadOrCreateCA(caKeyFile, caCertFile)
	if err != nil {
		return nil, fmt.Errorf("load or create CA key pair: %w", err)
	}

	dbLogger := cfg.Logger.Named("boltdb").Sugar()
	boltOpts := *bbolt.DefaultOptions
	boltOpts.Logger = &bolt.Logger{SugaredLogger: dbLogger}

	boltDB, err := bolt.OpenDatabase(dbPath, &boltOpts)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	scopeSvc := &scope.Scope{}

	reqLogService := reqlog.NewService(reqlog.Config{
		Scope:      scopeSvc,
		Repository: boltDB,
		Logger:     cfg.Logger.Named("reqlog").Sugar(),
	})

	interceptService := intercept.NewService(intercept.Config{
		Logger: cfg.Logger.Named("intercept").Sugar(),
	})

	senderService := sender.NewService(sender.Config{
		Repository:    boltDB,
		ReqLogService: reqLogService,
	})

	projService, err := proj.NewService(proj.Config{
		Repository:       boltDB,
		InterceptService: interceptService,
		ReqLogService:    reqLogService,
		SenderService:    senderService,
		Scope:            scopeSvc,
	})
	if err != nil {
		boltDB.Close()
		return nil, fmt.Errorf("create projects service: %w", err)
	}

	proxyHandler, err := proxy.NewProxy(proxy.Config{
		CACert: caCert,
		CAKey:  caKey,
		Logger: cfg.Logger.Named("proxy").Sugar(),
	})
	if err != nil {
		boltDB.Close()
		return nil, fmt.Errorf("create proxy: %w", err)
	}

	proxyHandler.UseRequestModifier(reqLogService.RequestModifier)
	proxyHandler.UseResponseModifier(reqLogService.ResponseModifier)
	proxyHandler.UseRequestModifier(interceptService.RequestModifier)
	proxyHandler.UseResponseModifier(interceptService.ResponseModifier)

	adminHandler := http.FileServer(http.FS(cfg.AdminFS))
	router := mux.NewRouter().SkipClean(true)
	adminRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		hostname, _ := os.Hostname()
		host, _, _ := net.SplitHostPort(req.Host)

		return strings.EqualFold(host, hostname) ||
			req.Host == "hetty.proxy" ||
			req.Host == fmt.Sprintf("%v:%v", "localhost", listenPort) ||
			req.Host == fmt.Sprintf("%v:%v", listenHost, listenPort) ||
			req.Method != http.MethodConnect && !strings.HasPrefix(req.RequestURI, "http://")
	}).Subrouter().StrictSlash(true)

	gqlEndpoint := "/api/graphql/"
	adminRouter.Path(gqlEndpoint).Handler(api.HTTPHandler(&api.Resolver{
		ProjectService:    projService,
		RequestLogService: reqLogService,
		InterceptService:  interceptService,
		SenderService:     senderService,
	}, gqlEndpoint))

	adminRouter.PathPrefix("").Handler(adminHandler)
	router.PathPrefix("").Handler(proxyHandler)

	httpServer := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){},
		ErrorLog:     zap.NewStdLog(cfg.Logger.Named("http")),
	}

	inst := &Instance{
		HTTPServer: httpServer,
		URL:        url,
		boltDB:     boltDB,
		logger:     mainLogger,
	}

	go func() {
		mainLogger.Info(fmt.Sprintf("Hetty (v%v) is running on %v ...", cfg.Version, cfg.Addr))
		mainLogger.Info(fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(32), "Get started at "+url))

		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			mainLogger.Fatal("HTTP server closed unexpected.", zap.Error(err))
		}
	}()

	return inst, nil
}

// Shutdown gracefully stops the HTTP server and closes the database.
func (i *Instance) Shutdown(ctx context.Context) error {
	i.logger.Info("Shutting down HTTP server.")

	err := i.HTTPServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	if err := i.boltDB.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}

	return nil
}

// WaitReady polls the admin URL until the server responds or the timeout elapses.
func WaitReady(url string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server not ready at %s within %v", url, timeout)
}
