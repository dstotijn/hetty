package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
)

// Proxy is used to forward HTTP requests.
type Proxy struct {
	rp httputil.ReverseProxy
}

// NewProxy returns a new Proxy.
func NewProxy() *Proxy {
	return &Proxy{
		rp: httputil.ReverseProxy{
			Director: func(r *http.Request) {
				log.Printf("Director handled URL: %v", r.URL)
			},
		},
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("ServeHTTP: Received request (host: %v, url: %v", r.Host, r.URL)

	if r.Method == http.MethodConnect {
		p.handleConnect(w, r)
		return
	}

	p.rp.ServeHTTP(w, r)
	log.Printf("ServeHTTP: Finished (host: %v, url: %v", r.Host, r.URL)
}

func (p *Proxy) handleConnect(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Printf("handleConnect: ResponseWriter is not a http.Hijacker (type: %T)", w)
		writeError(w, r, http.StatusServiceUnavailable)
		return
	}

	// destConn is the TCP connection to the destination web server of the
	// proxied HTTP request.
	destConn, err := net.Dial("tcp", r.Host)
	if err != nil {
		log.Printf("handleConnect: Connect to destination host failed: %v", err)
		writeError(w, r, http.StatusBadGateway)
		return
	}
	defer destConn.Close()

	w.WriteHeader(http.StatusOK)

	// clientConn is the TCP connection to the client.
	clientConn, _, err := hj.Hijack()
	if err != nil {
		log.Printf("handleConnect: Hijack failed: %v", err)
		writeError(w, r, http.StatusServiceUnavailable)
		return
	}
	defer clientConn.Close()

	errc := make(chan error, 1)
	go tunnelData(destConn, clientConn, errc)
	go tunnelData(clientConn, destConn, errc)
	<-errc
}

func tunnelData(dst, src io.ReadWriter, errc chan<- error) {
	_, err := io.Copy(dst, src)
	errc <- err
}

func writeError(w http.ResponseWriter, r *http.Request, code int) {
	http.Error(w, http.StatusText(code), code)
}
