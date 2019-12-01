package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/99designs/gqlgen/handler"
	"github.com/dstotijn/gurp/pkg/api"
	"github.com/dstotijn/gurp/pkg/proxy"
	"github.com/dstotijn/gurp/pkg/reqlog"

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

	reqLog := reqlog.NewRequestLog()

	p, err := proxy.NewProxy(caCert, tlsCA.PrivateKey)
	if err != nil {
		log.Fatalf("[FATAL] Could not create Proxy: %v", err)
	}

	p.UseRequestModifier(func(next proxy.RequestModifyFunc) proxy.RequestModifyFunc {
		return func(req *http.Request) {
			next(req)
			clone := req.Clone(req.Context())
			var body []byte
			if req.Body != nil {
				// TODO: Use io.LimitReader.
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					log.Printf("[ERROR] Could not read request body for logging: %v", err)
					return
				}
				req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
			reqLog.AddRequest(*clone, body)
		}
	})

	p.UseResponseModifier(func(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
		return func(res *http.Response) error {
			if err := next(res); err != nil {
				return err
			}
			clone := *res
			var body []byte
			if res.Body != nil {
				// TODO: Use io.LimitReader.
				var err error
				body, err = ioutil.ReadAll(res.Body)
				if err != nil {
					return fmt.Errorf("could not read response body: %v", err)
				}
				res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
			reqLog.AddResponse(clone, body)
			return nil
		}
	})

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
		return strings.EqualFold(host, hostname) || req.Host == "gurp.proxy"
	}).Subrouter()

	// GraphQL server.
	adminRouter.Path("/api/playground").Handler(handler.Playground("GraphQL Playground", "/api/graphql"))
	adminRouter.Path("/api/graphql").Handler(handler.GraphQL(api.NewExecutableSchema(api.Config{Resolvers: &api.Resolver{}})))

	// Admin interface.
	adminRouter.PathPrefix("").Handler(adminHandler)

	// Fallback (default) is the Proxy handler.
	router.PathPrefix("").Handler(p)

	s := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[FATAL] HTTP server closed: %v", err)
	}
}
