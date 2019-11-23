package main

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

var httpHandler = &httputil.ReverseProxy{
	Director: func(r *http.Request) {
		r.URL.Host = r.Host
		r.URL.Scheme = "http"
	},
	ErrorHandler: proxyErrorHandler,
}

var httpsHandler = &httputil.ReverseProxy{
	Director: func(r *http.Request) {
		r.URL.Host = r.Host
		r.URL.Scheme = "https"
	},
	ErrorHandler: proxyErrorHandler,
}

func proxyErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == context.Canceled {
		return
	}
	log.Printf("[ERROR]: Proxy error: %v", err)
	w.WriteHeader(http.StatusBadGateway)
}

// Proxy is used to forward HTTP requests.
type Proxy struct {
	certConfig *CertConfig
}

// NewProxy returns a new Proxy.
func NewProxy(ca *x509.Certificate, key crypto.PrivateKey) (*Proxy, error) {
	certConfig, err := NewCertConfig(ca, key)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		certConfig: certConfig,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}

	httpHandler.ServeHTTP(w, r)
}

// handleConnect hijacks the incoming HTTP request and sets up an HTTP tunnel.
// During the TLS handshake with the client, we use the proxy's CA config to
// create a certificate on-the-fly.
func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Printf("[ERROR] handleConnect: ResponseWriter is not a http.Hijacker (type: %T)", w)
		writeError(w, r, http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Printf("[ERROR] Hijacking client connection failed: %v", err)
		writeError(w, r, http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	// Secure connection to client.
	clientConn, err = p.clientTLSConn(clientConn)
	if err != nil {
		log.Printf("[ERROR] Securing client connection failed: %v", err)
		return
	}
	clientConnNotify := ConnNotify{clientConn, make(chan struct{})}

	l := &OnceAcceptListener{clientConnNotify.Conn}

	err = http.Serve(l, httpsHandler)
	if err != nil && err != ErrAlreadyAccepted {
		log.Printf("[ERROR] Serving HTTP request failed: %v", err)
	}
	<-clientConnNotify.closed
}

func (p *Proxy) clientTLSConn(conn net.Conn) (*tls.Conn, error) {
	tlsConfig := p.certConfig.TLSConfig()

	tlsConn := tls.Server(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("handshake error: %v", err)
	}

	return tlsConn, nil
}

func writeError(w http.ResponseWriter, r *http.Request, code int) {
	http.Error(w, http.StatusText(code), code)
}
