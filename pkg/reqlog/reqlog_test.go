package reqlog_test

import (
	"context"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid"

	badgerdb "github.com/dgraph-io/badger/v3"
	"github.com/dstotijn/hetty/pkg/db/badger"
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

//nolint:paralleltest
func TestRequestModifier(t *testing.T) {
	db, err := badger.OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	svc := reqlog.NewService(reqlog.Config{
		Repository: db,
		Scope:      &scope.Scope{},
	})
	svc.SetActiveProjectID(ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy))

	next := func(req *http.Request) {
		req.Body = io.NopCloser(strings.NewReader("modified body"))
	}
	reqModFn := svc.RequestModifier(next)
	req := httptest.NewRequest("GET", "https://example.com/", strings.NewReader("bar"))
	reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	req = req.WithContext(proxy.WithRequestID(req.Context(), reqID))

	reqModFn(req)

	t.Run("request log was stored in repository", func(t *testing.T) {
		exp := reqlog.RequestLog{
			ID:        reqID,
			ProjectID: svc.ActiveProjectID(),
			Method:    req.Method,
			URL:       req.URL,
			Proto:     req.Proto,
			Header:    req.Header,
			Body:      []byte("modified body"),
		}

		got, err := svc.FindRequestLogByID(context.Background(), reqID)
		if err != nil {
			t.Fatalf("failed to find request by id: %v", err)
		}

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("request log not equal (-exp, +got):\n%v", diff)
		}
	})
}

//nolint:paralleltest
func TestResponseModifier(t *testing.T) {
	db, err := badger.OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	svc := reqlog.NewService(reqlog.Config{
		Repository: db,
	})
	svc.SetActiveProjectID(ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy))

	next := func(res *http.Response) error {
		res.Body = io.NopCloser(strings.NewReader("modified body"))
		return nil
	}
	resModFn := svc.ResponseModifier(next)

	req := httptest.NewRequest("GET", "https://example.com/", strings.NewReader("bar"))
	reqLogID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	req = req.WithContext(context.WithValue(req.Context(), reqlog.ReqLogIDKey, reqLogID))

	err = db.StoreRequestLog(context.Background(), reqlog.RequestLog{
		ID: reqLogID,
	})
	if err != nil {
		t.Fatalf("failed to store request log: %v", err)
	}

	res := &http.Response{
		Request: req,
		Body:    io.NopCloser(strings.NewReader("bar")),
	}

	if err := resModFn(res); err != nil {
		t.Fatalf("unexpected error (expected: nil, got: %v)", err)
	}

	t.Run("request log was stored in repository", func(t *testing.T) {
		// Dirty (but simple) wait for other goroutine to finish calling repository.
		time.Sleep(10 * time.Millisecond)

		got, err := svc.FindRequestLogByID(context.Background(), reqLogID)
		if err != nil {
			t.Fatalf("failed to find request by id: %v", err)
		}

		t.Run("ran next modifier first, before calling repository", func(t *testing.T) {
			if exp := "modified body"; exp != string(got.Response.Body) {
				t.Fatalf("incorrect `ResponseLog.Body` value (expected: %v, got: %v)", exp, string(got.Response.Body))
			}
		})
	})
}
