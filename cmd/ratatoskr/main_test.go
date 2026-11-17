package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestPlaceholderHandler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "root serves the placeholder",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "<h1>Ratatoskr</h1>",
		},
		{
			name:       "head root",
			method:     http.MethodHead,
			path:       "/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "unknown path is not the placeholder",
			method:     http.MethodGet,
			path:       "/nope",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "non-GET is rejected",
			method:     http.MethodPost,
			path:       "/",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			newHandler().ServeHTTP(rec, httptest.NewRequest(tt.method, tt.path, nil))

			if got := rec.Code; got != tt.wantStatus {
				t.Errorf("status = %d, want %d", got, tt.wantStatus)
			}

			if !strings.Contains(rec.Body.String(), tt.wantBody) {
				t.Errorf("body %q does not contain %q", rec.Body.String(), tt.wantBody)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			h := rec.Header()
			if got, want := h.Get("Content-Type"), "text/html; charset=utf-8"; got != want {
				t.Errorf("Content-Type = %q, want %q", got, want)
			}

			if got, want := h.Get("X-Content-Type-Options"), "nosniff"; got != want {
				t.Errorf("X-Content-Type-Options = %q, want %q", got, want)
			}
		})
	}
}

func TestServePlaceholder(t *testing.T) {
	t.Parallel()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	serveErr := make(chan error, 1)
	go func() { serveErr <- serve(ctx, ln, options{}) }()

	res, err := http.Get("http://" + ln.Addr().String() + "/")
	if err != nil {
		t.Fatalf("get placeholder: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if closeErr := res.Body.Close(); closeErr != nil {
		t.Fatalf("close body: %v", closeErr)
	}

	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", res.StatusCode, http.StatusOK)
	}

	if !strings.Contains(string(body), "Ratatoskr") {
		t.Errorf("body %q does not contain %q", body, "Ratatoskr")
	}

	cancel()

	select {
	case err := <-serveErr:
		if err != nil {
			t.Fatalf("serve: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("serve did not return after the context was cancelled")
	}
}

func TestParseFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
		want options
	}{
		{
			name: "defaults",
			want: options{listen: "127.0.0.1:8080"},
		},
		{
			name: "listen and data dir",
			args: []string{"--listen", "127.0.0.1:0", "--data-dir", "./.dev-data"},
			want: options{listen: "127.0.0.1:0", dataDir: "./.dev-data"},
		},
		{
			name: "version",
			args: []string{"-version"},
			want: options{listen: "127.0.0.1:8080", showVersion: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseFlags(tt.args)
			if err != nil {
				t.Fatalf("parseFlags(%q): %v", tt.args, err)
			}

			if got != tt.want {
				t.Errorf("parseFlags(%q) = %+v, want %+v", tt.args, got, tt.want)
			}
		})
	}

	if _, err := parseFlags([]string{"--nope"}); err == nil {
		t.Error("parseFlags(--nope) = nil error, want an error")
	}
}
