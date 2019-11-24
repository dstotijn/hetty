package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"

	"github.com/dstotijn/gurp/proxy"
)

var (
	caCertFile = flag.String("cert", "", "CA certificate file path")
	caKeyFile  = flag.String("key", "", "CA private key file path")
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

	proxy, err := proxy.NewProxy(caCert, tlsCA.PrivateKey)
	if err != nil {
		log.Fatalf("[FATAL] Could not create Proxy: %v", err)
	}

	proxy.UseRequestModifier(func(req *http.Request) {
		log.Printf("[DEBUG] Incoming request: %v", req.URL)
	})

	proxy.UseResponseModifier(func(res *http.Response) error {
		log.Printf("[DEBUG] Downstream response: %v %v %v", res.Proto, res.StatusCode, http.StatusText(res.StatusCode))
		return nil
	})

	s := &http.Server{
		Addr:         ":8080",
		Handler:      proxy,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	err = s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("[FATAL] HTTP server closed: %v", err)
	}
}
