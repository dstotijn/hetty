package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io/ioutil"
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

func changeBody(res *http.Response, modifer func(body []byte) []byte) error {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Could not read response body: %v", err)
	}

	contentEncoding := res.Header.Get("Content-Encoding")

	if contentEncoding == "" {
		newBody := modifer(body)
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))
	}

	if contentEncoding == "gzip" {
		// TMP!
		//res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			return fmt.Errorf("Could not create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		body, err = ioutil.ReadAll(gzipReader)

		// TODO: Gzip this body
		newBody := modifer(body)
		res.Header.Set("Content-Encoding", "")
		res.Body = ioutil.NopCloser(bytes.NewBuffer(newBody))

		if err != nil {
			return fmt.Errorf("Could not read gzipped response body: %v", err)
		}
	}

	return nil
}

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

	_, err = projService.Open(context.TODO(), "New2")
	if err != nil {
		panic(err)
	}

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

	p.UseResponseModifier(func(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
		return func(res *http.Response) error {
			if err := next(res); err != nil {
				return err
			}

			if res.Request.URL.String() == "http://www.ue.wroc.pl/" {
				//res.Header["X-Proxy"] = []string{"Hello"}

				fmt.Println("Intercepting ue.wroc.pl request!")

				err := changeBody(res, func(b []byte) []byte {
					//fmt.Println(string(b))

					prefix := []byte("<p>Pozdro poćwicz</p>")

					return append(prefix, b...)
				})
				if err != nil {
					panic(err)
				}

				//res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}

			return nil
		}
	})

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
