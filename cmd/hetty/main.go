package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dstotijn/hetty/pkg/api"
	"github.com/dstotijn/hetty/pkg/db/sqlite"
	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
	"github.com/mitchellh/go-homedir"
)

var (
	caCertFile string
	caKeyFile  string
	projPath   string
	addr       string
	adminPath  string
)

func main() {
	flag.StringVar(&caCertFile, "cert", "~/.hetty/hetty_cert.pem", "CA certificate filepath. Creates a new CA certificate is file doesn't exist")
	flag.StringVar(&caKeyFile, "key", "~/.hetty/hetty_key.pem", "CA private key filepath. Creates a new CA private key if file doesn't exist")
	flag.StringVar(&projPath, "projects", "~/.hetty/projects", "Projects directory path")
	flag.StringVar(&addr, "addr", ":8080", "TCP address to listen on, in the form \"host:port\"")
	flag.StringVar(&adminPath, "adminPath", "", "File path to admin build")
	flag.Parse()

	// Expand `~` in filepaths.
	caCertFile, err := homedir.Expand(caCertFile)
	if err != nil {
		log.Fatalf("[FATAL] Could not parse CA certificate filepath: %v", err)
	}
	caKeyFile, err := homedir.Expand(caKeyFile)
	if err != nil {
		log.Fatalf("[FATAL] Could not parse CA private key filepath: %v", err)
	}
	projPath, err := homedir.Expand(projPath)
	if err != nil {
		log.Fatalf("[FATAL] Could not parse projects filepath: %v", err)
	}

	// Load existing CA certificate and key from disk, or generate and write
	// to disk if no files exist yet.
	caCert, caKey, err := proxy.LoadOrCreateCA(caKeyFile, caCertFile)
	if err != nil {
		log.Fatalf("[FATAL] Could not create/load CA key pair: %v", err)
	}

	db, err := sqlite.New(projPath)
	if err != nil {
		log.Fatalf("[FATAL] Could not initialize database client: %v", err)
	}

	projService, err := proj.NewService(db)
	if err != nil {
		log.Fatalf("[FATAL] Could not create new project service: %v", err)
	}
	defer projService.Close()

	scope := scope.New(db, projService)

	reqLogService := reqlog.NewService(reqlog.Config{
		Scope:          scope,
		ProjectService: projService,
		Repository:     db,
	})

	p, err := proxy.NewProxy(caCert, caKey)
	if err != nil {
		log.Fatalf("[FATAL] Could not create Proxy: %v", err)
	}

	p.UseRequestModifier(reqLogService.RequestModifier)
	p.UseResponseModifier(reqLogService.ResponseModifier)

	var adminHandler http.Handler
	if adminPath == "" {
		// Used for embedding with `rice`.
		box, err := rice.FindBox("../../admin/dist")
		if err != nil {
			log.Fatalf("[FATAL] Could not find embedded admin resources: %v", err)
		}
		adminHandler = http.FileServer(box.HTTPBox())
	} else {
		adminHandler = http.FileServer(http.Dir(adminPath))
	}

	router := mux.NewRouter().SkipClean(true)

	adminRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		hostname, _ := os.Hostname()
		host, _, _ := net.SplitHostPort(req.Host)
		return strings.EqualFold(host, hostname) || (req.Host == "hetty.proxy" || req.Host == "localhost:8080")
	}).Subrouter().StrictSlash(true)

	// GraphQL server.
	adminRouter.Path("/api/playground/").Handler(playground.Handler("GraphQL Playground", "/api/graphql/"))
	adminRouter.Path("/api/graphql/").Handler(handler.NewDefaultServer(api.NewExecutableSchema(api.Config{Resolvers: &api.Resolver{
		RequestLogService: reqLogService,
		ProjectService:    projService,
		ScopeService:      scope,
	}})))

	// Admin interface.
	adminRouter.PathPrefix("").Handler(adminHandler)

	// Fallback (default) is the Proxy handler.
	router.PathPrefix("").Handler(p)

	s := &http.Server{
		Addr:         addr,
		Handler:      router,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	log.Printf("[INFO] Running server on %v ...", addr)
	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[FATAL] HTTP server closed: %v", err)
	}
}
