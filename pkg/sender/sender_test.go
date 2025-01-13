package sender_test

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
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"

	"github.com/dstotijn/hetty/pkg/db/bolt"
	"github.com/dstotijn/hetty/pkg/proj"
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

		path := t.TempDir() + "bolt.db"
		boltDB, err := bbolt.Open(path, 0o600, nil)
		if err != nil {
			t.Fatalf("failed to open bolt database: %v", err)
		}
		defer boltDB.Close()

		db, err := bolt.DatabaseFromBoltDB(boltDB)
		if err != nil {
			t.Fatalf("failed to create database: %v", err)
		}
		defer db.Close()

		svc := sender.NewService(sender.Config{
			Repository: db,
		})

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		err = db.UpsertProject(context.Background(), proj.Project{
			ID:       projectID,
			Name:     "foobar",
			Settings: proj.Settings{},
		})
		if err != nil {
			t.Fatalf("unexpected error upserting project: %v", err)
		}

		svc.SetActiveProjectID(projectID)
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

		diff := cmp.Diff(exp, got, cmpopts.IgnoreFields(sender.Request{}, "ID"))
		if diff != "" {
			t.Fatalf("request not equal (-exp, +got):\n%v", diff)
		}

		got, err = db.FindSenderRequestByID(context.Background(), projectID, got.ID)
		if err != nil {
			t.Fatalf("failed to find request by ID: %v", err)
		}

		diff = cmp.Diff(exp, got, cmpopts.IgnoreFields(sender.Request{}, "ID"))
		if diff != "" {
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

		path := t.TempDir() + "bolt.db"
		boltDB, err := bbolt.Open(path, 0o600, nil)
		if err != nil {
			t.Fatalf("failed to open bolt database: %v", err)
		}
		defer boltDB.Close()

		db, err := bolt.DatabaseFromBoltDB(boltDB)
		if err != nil {
			t.Fatalf("failed to create database: %v", err)
		}
		defer db.Close()

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
		err = db.UpsertProject(context.Background(), proj.Project{
			ID: projectID,
		})
		if err != nil {
			t.Fatalf("unexpected error upserting project: %v", err)
		}

		reqLog := reqlog.RequestLog{
			ID:        reqLogID,
			ProjectID: projectID,
			URL:       exampleURL,
			Method:    http.MethodPost,
			Proto:     "HTTP/1.1",
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		}

		if err := db.StoreRequestLog(context.Background(), reqLog); err != nil {
			t.Fatalf("failed to store request log: %v", err)
		}

		svc := sender.NewService(sender.Config{
			ReqLogService: reqlog.NewService(reqlog.Config{
				ActiveProjectID: projectID,
				Repository:      db,
			}),
			Repository: db,
		})

		svc.SetActiveProjectID(projectID)

		exp := sender.Request{
			SourceRequestLogID: reqLogID,
			ProjectID:          projectID,
			URL:                exampleURL,
			Method:             http.MethodPost,
			Proto:              sender.HTTPProto20,
			Header: http.Header{
				"X-Foo": []string{"bar"},
			},
			Body: []byte("foobar"),
		}

		got, err := svc.CloneFromRequestLog(context.Background(), reqLogID)
		if err != nil {
			t.Fatalf("unexpected error cloning from request log: %v", err)
		}

		diff := cmp.Diff(exp, got, cmpopts.IgnoreFields(sender.Request{}, "ID"))
		if diff != "" {
			t.Fatalf("request not equal (-exp, +got):\n%v", diff)
		}
	})
}

func TestSendRequest(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "bolt.db"
	boltDB, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		t.Fatalf("failed to open bolt database: %v", err)
	}
	defer boltDB.Close()

	db, err := bolt.DatabaseFromBoltDB(boltDB)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	date := time.Now().Format(http.TimeFormat)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Foobar", "baz")
		w.Header().Set("Date", date)
		fmt.Fprint(w, "baz")
	}))
	defer ts.Close()

	tsURL, _ := url.Parse(ts.URL)

	projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	err = db.UpsertProject(context.Background(), proj.Project{
		ID:       projectID,
		Settings: proj.Settings{},
	})
	if err != nil {
		t.Fatalf("unexpected error upserting project: %v", err)
	}

	reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	req := sender.Request{
		ID:        reqID,
		ProjectID: projectID,
		URL:       tsURL,
		Method:    http.MethodPost,
		Proto:     "HTTP/1.1",
		Header: http.Header{
			"X-Foo": []string{"bar"},
		},
		Body: []byte("foobar"),
	}

	if err := db.StoreSenderRequest(context.Background(), req); err != nil {
		t.Fatalf("failed to store request: %v", err)
	}

	svc := sender.NewService(sender.Config{
		ReqLogService: reqlog.NewService(reqlog.Config{
			Repository: db,
		}),
		Repository: db,
	})
	svc.SetActiveProjectID(projectID)

	exp := &reqlog.ResponseLog{
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

	diff := cmp.Diff(exp, got.Response, cmpopts.IgnoreFields(sender.Request{}, "ID"))
	if diff != "" {
		t.Fatalf("request not equal (-exp, +got):\n%v", diff)
	}
}
