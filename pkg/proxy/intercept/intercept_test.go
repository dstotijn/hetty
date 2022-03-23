package intercept_test

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/oklog/ulid"
	"go.uber.org/zap"

	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/proxy/intercept"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

func TestRequestModifier(t *testing.T) {
	t.Parallel()

	t.Run("modify request that's not found", func(t *testing.T) {
		t.Parallel()

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  true,
			ResponsesEnabled: false,
		})

		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

		err := svc.ModifyRequest(reqID, nil, nil)
		if !errors.Is(err, intercept.ErrRequestNotFound) {
			t.Fatalf("expected `intercept.ErrRequestNotFound`, got: %v", err)
		}
	})

	t.Run("modify request that's done", func(t *testing.T) {
		t.Parallel()

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  true,
			ResponsesEnabled: false,
		})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req := httptest.NewRequest("GET", "https://example.com/foo", nil)
		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		*req = *req.WithContext(ctx)
		*req = *req.WithContext(proxy.WithRequestID(req.Context(), reqID))

		next := func(req *http.Request) {}
		go svc.RequestModifier(next)(req)

		// Wait shortly, to allow the req modifier goroutine to add `req` to the
		// array of intercepted reqs.
		time.Sleep(10 * time.Millisecond)
		cancel()

		modReq := req.Clone(req.Context())
		modReq.Header.Set("X-Foo", "bar")

		err := svc.ModifyRequest(reqID, modReq, nil)
		if !errors.Is(err, intercept.ErrRequestDone) {
			t.Fatalf("expected `intercept.ErrRequestDone`, got: %v", err)
		}
	})

	t.Run("modify intercepted request", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "https://example.com/foo", nil)
		req.Header.Set("X-Foo", "foo")

		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		*req = *req.WithContext(proxy.WithRequestID(req.Context(), reqID))

		modReq := req.Clone(context.Background())
		modReq.Header.Set("X-Foo", "bar")

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  true,
			ResponsesEnabled: false,
		})

		var got *http.Request

		next := func(req *http.Request) {
			got = req.Clone(context.Background())
		}

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			svc.RequestModifier(next)(req)
			wg.Done()
		}()

		// Wait shortly, to allow the req modifier goroutine to add `req` to the
		// array of intercepted reqs.
		time.Sleep(10 * time.Millisecond)

		err := svc.ModifyRequest(reqID, modReq, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wg.Wait()

		if got == nil {
			t.Fatal("expected `got` not to be nil")
		}

		if exp := "bar"; exp != got.Header.Get("X-Foo") {
			t.Fatalf("incorrect modified request header value (expected: %v, got: %v)", exp, got.Header.Get("X-Foo"))
		}
	})
}

func TestResponseModifier(t *testing.T) {
	t.Parallel()

	t.Run("modify response that's not found", func(t *testing.T) {
		t.Parallel()

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  false,
			ResponsesEnabled: true,
		})

		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

		err := svc.ModifyResponse(reqID, nil)
		if !errors.Is(err, intercept.ErrRequestNotFound) {
			t.Fatalf("expected `intercept.ErrRequestNotFound`, got: %v", err)
		}
	})

	t.Run("modify response of request that's done", func(t *testing.T) {
		t.Parallel()

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  false,
			ResponsesEnabled: true,
		})

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		req := httptest.NewRequest("GET", "https://example.com/foo", nil)
		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		*req = *req.WithContext(ctx)
		*req = *req.WithContext(proxy.WithRequestID(req.Context(), reqID))

		res := &http.Response{
			Request: req,
			Header:  make(http.Header),
		}
		res.Header.Add("X-Foo", "foo")

		var modErr error
		var wg sync.WaitGroup
		wg.Add(1)

		next := func(res *http.Response) error { return nil }
		go func() {
			defer wg.Done()
			modErr = svc.ResponseModifier(next)(res)
		}()

		// Wait shortly, to allow the res modifier goroutine to add `res` to the
		// array of intercepted responses.
		time.Sleep(10 * time.Millisecond)
		cancel()

		modRes := *res
		modRes.Header = make(http.Header)
		modRes.Header.Set("X-Foo", "bar")

		err := svc.ModifyResponse(reqID, &modRes)
		if !errors.Is(err, intercept.ErrRequestDone) {
			t.Fatalf("expected `intercept.ErrRequestDone`, got: %v", err)
		}

		wg.Wait()

		if !errors.Is(modErr, context.Canceled) {
			t.Fatalf("expected `context.Canceled`, got: %v", modErr)
		}
	})

	t.Run("modify intercepted response", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest("GET", "https://example.com/foo", nil)
		req.Header.Set("X-Foo", "foo")

		reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		*req = *req.WithContext(proxy.WithRequestID(req.Context(), reqID))

		res := &http.Response{
			Request: req,
			Header:  make(http.Header),
		}
		res.Header.Add("X-Foo", "foo")

		modRes := *res
		modRes.Header = make(http.Header)
		modRes.Header.Set("X-Foo", "bar")

		logger, _ := zap.NewDevelopment()
		svc := intercept.NewService(intercept.Config{
			Logger:           logger.Sugar(),
			RequestsEnabled:  false,
			ResponsesEnabled: true,
		})

		var gotHeader string

		var next proxy.ResponseModifyFunc = func(res *http.Response) error {
			gotHeader = res.Header.Get("X-Foo")
			return nil
		}

		var modErr error
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			modErr = svc.ResponseModifier(next)(res)
			wg.Done()
		}()

		// Wait shortly, to allow the res modifier goroutine to add `req` to the
		// array of intercepted reqs.
		time.Sleep(10 * time.Millisecond)

		err := svc.ModifyResponse(reqID, &modRes)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		wg.Wait()

		if modErr != nil {
			t.Fatalf("unexpected error: %v", modErr)
		}

		if exp := "bar"; exp != gotHeader {
			t.Fatalf("incorrect modified request header value (expected: %v, got: %v)", exp, gotHeader)
		}
	})
}
