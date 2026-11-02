package sender

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

type HTTPTransport struct{}

type protoCtxKey struct{}

const (
	HTTPProto10 = "HTTP/1.0"
	HTTPProto11 = "HTTP/1.1"
	HTTPProto20 = "HTTP/2.0"
)

// defaultHTTPTransport is used for outbound requests in MITM mode. The sender
// intentionally skips certificate verification because it may connect to HTTPS
// targets using IP literals that do not include SAN entries.
var defaultHTTPTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},
}

// h1OnlyTransport mimics `http.DefaultTransport`, but with HTTP/2 disabled.
var h1OnlyTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true,
	},

	// Disable HTTP/2.
	TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
}

// RountTrip implements http.RoundTripper. Based on a context value on the
// HTTP request, it switches between using `http.DefaultTransport` (which attempts
// HTTP/2) and a HTTP/1.1 only transport that's based off `http.DefaultTransport`.
func (t *HTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	proto, ok := req.Context().Value(protoCtxKey{}).(string)

	if ok && proto == HTTPProto10 || proto == HTTPProto11 {
		return h1OnlyTransport.RoundTrip(req)
	}

	return defaultHTTPTransport.RoundTrip(req)
}

func isValidProto(proto string) bool {
	return proto == HTTPProto10 || proto == HTTPProto11 || proto == HTTPProto20
}
