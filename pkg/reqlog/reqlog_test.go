package reqlog_test

//go:generate go run github.com/matryer/moq -out repo_mock_test.go -pkg reqlog_test . Repository:RepoMock

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

	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

//nolint:paralleltest
func TestRequestModifier(t *testing.T) {
	repoMock := &RepoMock{
		StoreRequestLogFunc: func(_ context.Context, _ reqlog.RequestLog) error {
			return nil
		},
	}
	svc := reqlog.NewService(reqlog.Config{
		Repository: repoMock,
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
		gotCount := len(repoMock.StoreRequestLogCalls())
		if expCount := 1; expCount != gotCount {
			t.Fatalf("incorrect `proj.Service.AddRequestLog` calls (expected: %v, got: %v)", expCount, gotCount)
		}

		exp := reqlog.RequestLog{
			ID:        ulid.ULID{}, // Empty value
			ProjectID: svc.ActiveProjectID(),
			Method:    req.Method,
			URL:       req.URL,
			Proto:     req.Proto,
			Header:    req.Header,
			Body:      []byte("modified body"),
		}
		got := repoMock.StoreRequestLogCalls()[0].ReqLog
		got.ID = ulid.ULID{} // Override to empty value so we can compare against expected value.

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("request log not equal (-exp, +got):\n%v", diff)
		}
	})
}

//nolint:paralleltest
func TestResponseModifier(t *testing.T) {
	repoMock := &RepoMock{
		StoreResponseLogFunc: func(_ context.Context, _ ulid.ULID, _ reqlog.ResponseLog) error {
			return nil
		},
	}
	svc := reqlog.NewService(reqlog.Config{
		Repository: repoMock,
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
		got := len(repoMock.StoreResponseLogCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.AddResponseLog` calls (expected: %v, got: %v)", exp, got)
		}

		t.Run("ran next modifier first, before calling repository", func(t *testing.T) {
			got := repoMock.StoreResponseLogCalls()[0].ResLog.Body
			if exp := "modified body"; exp != string(got) {
				t.Fatalf("incorrect `ResponseLog.Body` value (expected: %v, got: %v)", exp, string(got))
			}
		})

		t.Run("called repository with request log id", func(t *testing.T) {
			got := repoMock.StoreResponseLogCalls()[0].ReqLogID
			if exp := reqLogID; exp.Compare(got) != 0 {
				t.Fatalf("incorrect `reqLogID` argument for `Repository.AddResponseLogCalls` (expected: %v, got: %v)",
					exp.String(), got.String())
			}
		})
	})
}
