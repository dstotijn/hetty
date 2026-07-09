package proxy

import (
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/log"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

type contextKey int

const reqIDKey contextKey = 0

// Proxy implements http.Handler and offers MITM behaviour for modifying
// HTTP requests and responses.
type Proxy struct {
	certConfig *CertConfig
	handler    http.Handler
	logger     log.Logger

	// TODO: Add mutex for modifier funcs.
	reqModifiers []RequestModifyMiddleware
	resModifiers []ResponseModifyMiddleware
}

type Config struct {
	CACert *x509.Certificate
	CAKey  crypto.PrivateKey
	Logger log.Logger
}

// NewProxy returns a new Proxy.
func NewProxy(cfg Config) (*Proxy, error) {
	certConfig, err := NewCertConfig(cfg.CACert, cfg.CAKey)
	if err != nil {
		return nil, err
	}

	p := &Proxy{
		certConfig:   certConfig,
		reqModifiers: make([]RequestModifyMiddleware, 0),
		resModifiers: make([]ResponseModifyMiddleware, 0),
		logger:       cfg.Logger,
	}

	if p.logger == nil {
		p.logger = log.NewNopLogger()
	}

	transport := &http.Transport{
		// Values taken from `http.DefaultTransport`.
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

		// Non-default transport values.
		DisableCompression: true,
	}

	p.handler = &httputil.ReverseProxy{
		Transport:      transport,
		Director:       p.modifyRequest,
		ModifyResponse: p.modifyResponse,
		ErrorHandler:   p.errorHandler,
	}

	return p, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		p.handleConnect(w)
		return
	}

	reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	ctx := context.WithValue(r.Context(), reqIDKey, reqID)
	*r = *r.WithContext(ctx)

	p.handler.ServeHTTP(w, r)
}

func (p *Proxy) UseRequestModifier(fn ...RequestModifyMiddleware) {
	p.reqModifiers = append(p.reqModifiers, fn...)
}

func (p *Proxy) UseResponseModifier(fn ...ResponseModifyMiddleware) {
	p.resModifiers = append(p.resModifiers, fn...)
}

func (p *Proxy) modifyRequest(r *http.Request) {
	// Fix r.URL for HTTPS requests after CONNECT.
	if r.URL.Scheme == "" {
		r.URL.Host = r.Host
		r.URL.Scheme = "https"
	}

	// Setting `X-Forwarded-For` to `nil` ensures that http.ReverseProxy doesn't
	// set this header.
	r.Header["X-Forwarded-For"] = nil

	// Strip unsupported encodings.
	if acceptEncs := r.Header.Get("Accept-Encoding"); acceptEncs != "" {
		directives := strings.Split(acceptEncs, ",")
		updated := make([]string, 0, len(directives))

		for _, directive := range directives {
			stripped := strings.TrimSpace(directive)
			if strings.HasPrefix(stripped, "*") || strings.HasPrefix(stripped, "gzip") {
				updated = append(updated, stripped)
			}
		}

		if len(updated) == 0 {
			r.Header.Del("Accept-Encoding")
		} else {
			r.Header.Set("Accept-Encoding", strings.Join(updated, ", "))
		}
	}

	fn := nopReqModifier

	for i := len(p.reqModifiers) - 1; i >= 0; i-- {
		fn = p.reqModifiers[i](fn)
	}

	fn(r)
}

func (p *Proxy) modifyResponse(res *http.Response) error {
	fn := nopResModifier

	// TODO: Make decompressing gzip formatted response bodies a configurable project setting.
	if err := gunzipResponseBody(res); err != nil {
		return fmt.Errorf("proxy: failed to gunzip response body: %w", err)
	}

	for i := len(p.resModifiers) - 1; i >= 0; i-- {
		fn = p.resModifiers[i](fn)
	}

	return fn(res)
}

func WithRequestID(ctx context.Context, id ulid.ULID) context.Context {
	return context.WithValue(ctx, reqIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (ulid.ULID, bool) {
	id, ok := ctx.Value(reqIDKey).(ulid.ULID)
	return id, ok
}

// handleConnect hijacks the incoming HTTP request and sets up an HTTP tunnel.
// During the TLS handshake with the client, we use the proxy's CA config to
// create a certificate on-the-fly.
func (p *Proxy) handleConnect(w http.ResponseWriter) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		p.logger.Errorw("ResponseWriter is not a http.Hijacker.",
			"type", fmt.Sprintf("%T", w))
		writeError(w, http.StatusServiceUnavailable)

		return
	}

	w.WriteHeader(http.StatusOK)

	clientConn, _, err := hj.Hijack()
	if err != nil {
		p.logger.Errorw("Hijacking client connection failed.",
			"error", err)
		writeError(w, http.StatusServiceUnavailable)

		return
	}
	defer clientConn.Close()

	// Secure connection to client.
	tlsConn, err := p.clientTLSConn(clientConn)
	if err != nil {
		p.logger.Errorw("Securing client connection failed.",
			"error", err,
			"remoteAddr", clientConn.RemoteAddr().String())

		return
	}

	clientConnNotify := &ConnNotify{tlsConn, make(chan struct{}, 1)}
	l := &OnceAcceptListener{clientConnNotify}

	err = http.Serve(l, p)
	if err != nil && !errors.Is(err, ErrAlreadyAccepted) {
		p.logger.Errorw("Serving HTTP request failed.",
			"error", err)
	}

	<-clientConnNotify.closed
}

func (p *Proxy) clientTLSConn(conn net.Conn) (*tls.Conn, error) {
	tlsConfig := p.certConfig.TLSConfig()

	tlsConn := tls.Server(conn, tlsConfig)
	if err := tlsConn.Handshake(); err != nil {
		tlsConn.Close()
		return nil, fmt.Errorf("handshake error: %w", err)
	}

	return tlsConn, nil
}

func (p *Proxy) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case !errors.Is(err, context.Canceled):
		p.logger.Errorw("Failed to proxy request.",
			"error", err)
	case errors.Is(err, context.Canceled):
		p.logger.Debugw("Proxy request was cancelled.")
	}

	w.WriteHeader(http.StatusBadGateway)
}

func writeError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
