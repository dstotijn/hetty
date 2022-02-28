package sender_test

//go:generate go run github.com/matryer/moq -out reqlog_mock_test.go -pkg sender_test ../reqlog Service:ReqLogServiceMock
//go:generate go run github.com/matryer/moq -out repo_mock_test.go -pkg sender_test . Repository:RepoMock

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/sender"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

var exampleURL = func() *url.URL {
	u, err := url.Parse("https://example.com/foobar")
	if err != nil {
		panic(err)
	}

	return u
}()

func TestStoreRequest(t *testing.T) {
	t.Parallel()

	t.Run("without active project", func(t *testing.T) {
		t.Parallel()

		svc := sender.NewService(sender.Config{})

		_, err := svc.CreateOrUpdateRequest(context.Background(), sender.Request{
			URL:    exampleURL,
			Method: http.MethodPost,
			Body:   []byte("foobar"),
		})
		if !errors.Is(err, sender.ErrProjectIDMustBeSet) {
			t.Fatalf("expected `sender.ErrProjectIDMustBeSet`, got: %v", err)
		}
	})

	t.Run("with active project", func(t *testing.T) {
		t.Parallel()

		repoMock := &RepoMock{
			StoreSenderRequestFunc: func(ctx context.Context, req sender.Request) error {
				return nil
			},
		}
		svc := sender.NewService(sender.Config{
			Repository: repoMock,
		})

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		svc.SetActiveProjectID(projectID)

		exp := sender.Request{
			ProjectID: projectID,
			URL:       exampleURL,
			Method:    http.MethodPost,
			Proto:     "HTTP/1.1",
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		}

		got, err := svc.CreateOrUpdateRequest(context.Background(), sender.Request{
			URL:    exampleURL,
			Method: http.MethodPost,
			Proto:  "HTTP/1.1",
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		})
		if err != nil {
			t.Fatalf("unexpected error storing request: %v", err)
		}

		if got.ID.Compare(ulid.ULID{}) == 0 {
			t.Fatal("expected request ID to be non-empty value")
		}

		if len(repoMock.StoreSenderRequestCalls()) != 1 {
			t.Fatal("expected `svc.repo.StoreSenderRequest()` to have been called 1 time")
		}

		if diff := cmp.Diff(got, repoMock.StoreSenderRequestCalls()[0].Req); diff != "" {
			t.Fatalf("repo call arg not equal (-exp, +got):\n%v", diff)
		}

		// Reset ID to make comparison with expected easier.
		got.ID = ulid.ULID{}

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("request not equal (-exp, +got):\n%v", diff)
		}
	})
}

func TestCloneFromRequestLog(t *testing.T) {
	t.Parallel()

	reqLogID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

	t.Run("without active project", func(t *testing.T) {
		t.Parallel()

		svc := sender.NewService(sender.Config{})

		_, err := svc.CloneFromRequestLog(context.Background(), reqLogID)
		if !errors.Is(err, sender.ErrProjectIDMustBeSet) {
			t.Fatalf("expected `sender.ErrProjectIDMustBeSet`, got: %v", err)
		}
	})

	t.Run("with active project", func(t *testing.T) {
		t.Parallel()

		reqLog := reqlog.RequestLog{
			ID:        reqLogID,
			ProjectID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
			URL:       exampleURL,
			Method:    http.MethodPost,
			Proto:     "HTTP/1.1",
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		}

		reqLogMock := &ReqLogServiceMock{
			FindRequestLogByIDFunc: func(ctx context.Context, id ulid.ULID) (reqlog.RequestLog, error) {
				return reqLog, nil
			},
		}
		repoMock := &RepoMock{
			StoreSenderRequestFunc: func(ctx context.Context, req sender.Request) error {
				return nil
			},
		}
		svc := sender.NewService(sender.Config{
			ReqLogService: reqLogMock,
			Repository:    repoMock,
		})

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		svc.SetActiveProjectID(projectID)

		exp := sender.Request{
			SourceRequestLogID: reqLogID,
			ProjectID:          projectID,
			URL:                exampleURL,
			Method:             http.MethodPost,
			Proto:              sender.HTTPProto2,
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		}

		got, err := svc.CloneFromRequestLog(context.Background(), reqLogID)
		if err != nil {
			t.Fatalf("unexpected error cloning from request log: %v", err)
		}

		if len(reqLogMock.FindRequestLogByIDCalls()) != 1 {
			t.Fatal("expected `svc.reqLogSvc.FindRequestLogByID()` to have been called 1 time")
		}

		if got := reqLogMock.FindRequestLogByIDCalls()[0].ID; reqLogID.Compare(got) != 0 {
			t.Fatalf("reqlog service call arg `id` not equal (expected: %q, got: %q)", reqLogID, got)
		}

		if got.ID.Compare(ulid.ULID{}) == 0 {
			t.Fatal("expected request ID to be non-empty value")
		}

		if len(repoMock.StoreSenderRequestCalls()) != 1 {
			t.Fatal("expected `svc.repo.StoreSenderRequest()` to have been called 1 time")
		}

		if diff := cmp.Diff(got, repoMock.StoreSenderRequestCalls()[0].Req); diff != "" {
			t.Fatalf("repo call arg not equal (-exp, +got):\n%v", diff)
		}

		// Reset ID to make comparison with expected easier.
		got.ID = ulid.ULID{}

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("request not equal (-exp, +got):\n%v", diff)
		}
	})
}

func TestSendRequest(t *testing.T) {
	t.Parallel()

	date := time.Now().Format(http.TimeFormat)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Foobar", "baz")
		w.Header().Set("Date", date)
		fmt.Fprint(w, "baz")
	}))
	defer ts.Close()

	tsURL, _ := url.Parse(ts.URL)

	reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	req := sender.Request{
		ID:        reqID,
		ProjectID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
		URL:       tsURL,
		Method:    http.MethodPost,
		Proto:     "HTTP/1.1",
		Header: http.Header{
			"X-Foo": []string{"bar"},
		},
		Body: []byte("foobar"),
	}

	repoMock := &RepoMock{
		FindSenderRequestByIDFunc: func(ctx context.Context, id ulid.ULID) (sender.Request, error) {
			return req, nil
		},
		StoreResponseLogFunc: func(ctx context.Context, reqLogID ulid.ULID, resLog reqlog.ResponseLog) error {
			return nil
		},
	}
	svc := sender.NewService(sender.Config{
		Repository: repoMock,
	})

	exp := reqlog.ResponseLog{
		Proto:      "HTTP/1.1",
		StatusCode: http.StatusOK,
		Status:     "200 OK",
		Header: http.Header{
			"Content-Length": []string{"3"},
			"Content-Type":   []string{"text/plain; charset=utf-8"},
			"Date":           []string{date},
			"Foobar":         []string{"baz"},
		},
		Body: []byte("baz"),
	}

	got, err := svc.SendRequest(context.Background(), reqID)
	if err != nil {
		t.Fatalf("unexpected error sending request: %v", err)
	}

	if len(repoMock.FindSenderRequestByIDCalls()) != 1 {
		t.Fatal("expected `svc.repo.FindSenderRequestByID()` to have been called 1 time")
	}

	if diff := cmp.Diff(reqID, repoMock.FindSenderRequestByIDCalls()[0].ID); diff != "" {
		t.Fatalf("call arg `id` for `svc.repo.FindSenderRequestByID()` not equal (-exp, +got):\n%v", diff)
	}

	if len(repoMock.StoreResponseLogCalls()) != 1 {
		t.Fatal("expected `svc.repo.StoreResponseLog()` to have been called 1 time")
	}

	if diff := cmp.Diff(reqID, repoMock.StoreResponseLogCalls()[0].ReqLogID); diff != "" {
		t.Fatalf("call arg `reqLogID` for `svc.repo.StoreResponseLog()` not equal (-exp, +got):\n%v", diff)
	}

	if diff := cmp.Diff(exp, repoMock.StoreResponseLogCalls()[0].ResLog); diff != "" {
		t.Fatalf("call arg `resLog` for `svc.repo.StoreResponseLog()` not equal (-exp, +got):\n%v", diff)
	}

	if diff := cmp.Diff(repoMock.StoreResponseLogCalls()[0].ResLog, *got.Response); diff != "" {
		t.Fatalf("returned response log value and persisted value not equal (-exp, +got):\n%v", diff)
	}
}
