package main

import (
	"crypto/tls"
	"log"
	"net/http"
)

func main() {
	proxy := NewProxy()

	s := &http.Server{
		Addr:         ":8080",
		Handler:      proxy,
		TLSNextProto: map[string]func(*http.Server, *tls.Conn, http.Handler){}, // Disable HTTP/2
	}

	err := s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server closed: %v", err)
	}
}
