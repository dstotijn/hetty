package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/dstotijn/gurp/pkg/api"
	"github.com/dstotijn/gurp/pkg/proxy"
	"github.com/dstotijn/gurp/pkg/reqlog"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/mux"
)

var (
	caCertFile = flag.String("cert", "", "CA certificate file path")
	caKeyFile  = flag.String("key", "", "CA private key file path")
	dev        = flag.Bool("dev", false, "Run in development mode")
	adminPath  = flag.String("adminPath", "", "File path to admin build")
)

func main() {
	flag.Parse()

	tlsCA, err := tls.LoadX509KeyPair(*caCertFile, *caKeyFile)
	if err != nil {
		log.Fatalf("[FATAL] Could not load CA key pair: %v", err)
	}

	caCert, err := x509.ParseCertificate(tlsCA.Certificate[0])
	if err != nil {
		log.Fatalf("[FATAL] Could not parse CA: %v", err)
	}

	reqLogService := reqlog.NewService()

	p, err := proxy.NewProxy(caCert, tlsCA.PrivateKey)
	if err != nil {
		log.Fatalf("[FATAL] Could not create Proxy: %v", err)
	}

	p.UseRequestModifier(reqLogService.RequestModifier)
	p.UseResponseModifier(reqLogService.ResponseModifier)

	var adminHandler http.Handler

	if *dev {
		adminURL, err := url.Parse("http://localhost:3000")
		if err != nil {
			log.Fatalf("[FATAL] Invalid admin URL: %v", err)
		}
		adminHandler = httputil.NewSingleHostReverseProxy(adminURL)
	} else {
		if *adminPath == "" {
			log.Fatal("[FATAL] `adminPath` must be set")
		}
		adminHandler = http.FileServer(http.Dir(*adminPath))
	}

	router := mux.NewRouter().SkipClean(true)

	adminRouter := router.MatcherFunc(func(req *http.Request, match *mux.RouteMatch) bool {
		hostname, _ := os.Hostname()
		host, _, _ := net.SplitHostPort(req.Host)
		return strings.EqualFold(host, hostname) || (req.Host == "gurp.proxy" || req.Host == "localhost:8080")
	}).Subrouter()

	// GraphQL server.
	adminRouter.Path("/api/playground").Handler(playground.Handler("GraphQL Playground", "/api/graphql"))
	adminRouter.Path("/api/graphql").Handler(handler.NewDefaultServer(api.NewExecutableSchema(api.Config{Resolvers: &api.Resolver{
		RequestLogService: &reqLogService,
	}})))

	// Admin interface.
	adminRouter.PathPrefix("").Handler(adminHandler)

	// Fallback (default) is the Proxy handler.
	router.PathPrefix("").Handler(p)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	log.Println("[INFO] Running server on :8080 ...")
	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[FATAL] HTTP server closed: %v", err)
	}
}
