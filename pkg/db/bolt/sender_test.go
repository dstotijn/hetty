package bolt_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"

	"github.com/dstotijn/hetty/pkg/db/bolt"
	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/sender"
)

var exampleURL = func() *url.URL {
	u, err := url.Parse("https://example.com/foobar")
	if err != nil {
		panic(err)
	}

	return u
}()

func TestFindRequestByID(t *testing.T) {
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
	reqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

	err = db.UpsertProject(context.Background(), proj.Project{
		ID: projectID,
	})
	if err != nil {
		t.Fatalf("unexpected error upserting project: %v", err)
	}

	// See: https://go.dev/blog/subtests#cleaning-up-after-a-group-of-parallel-tests
	t.Run("group", func(t *testing.T) {
		t.Run("sender request not found", func(t *testing.T) {
			t.Parallel()

			_, err := db.FindSenderRequestByID(context.Background(), projectID, reqID)
			if !errors.Is(err, sender.ErrRequestNotFound) {
				t.Fatalf("expected `sender.ErrRequestNotFound`, got: %v", err)
			}
		})

		t.Run("sender request found", func(t *testing.T) {
			t.Parallel()

			exp := sender.Request{
				ID:                 ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				ProjectID:          projectID,
				SourceRequestLogID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),

				URL:    exampleURL,
				Method: http.MethodGet,
				Proto:  sender.HTTPProto20,
				Header: http.Header{
					"X-Foo": []string{"bar"},
				},
				Body: []byte("foo"),
				Response: &reqlog.ResponseLog{
					Proto:      "HTTP/2.0",
					Status:     "200 OK",
					StatusCode: 200,
					Header: http.Header{
						"X-Yolo": []string{"swag"},
					},
					Body: []byte("bar"),
				},
			}

			err := db.StoreSenderRequest(context.Background(), exp)
			if err != nil {
				t.Fatalf("unexpected error (expected: nil, got: %v)", err)
			}

			got, err := db.FindSenderRequestByID(context.Background(), exp.ProjectID, exp.ID)
			if err != nil {
				t.Fatalf("unexpected error (expected: nil, got: %v)", err)
			}

			if diff := cmp.Diff(exp, got); diff != "" {
				t.Fatalf("sender request not equal (-exp, +got):\n%v", diff)
			}
		})
	})
}

func TestFindSenderRequests(t *testing.T) {
	t.Parallel()

	t.Run("without project ID in filter", func(t *testing.T) {
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

		filter := sender.FindRequestsFilter{}

		_, err = db.FindSenderRequests(context.Background(), filter, nil)
		if !errors.Is(err, sender.ErrProjectIDMustBeSet) {
			t.Fatalf("expected `sender.ErrProjectIDMustBeSet`, got: %v", err)
		}
	})

	t.Run("returns sender requests and related response logs", func(t *testing.T) {
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
			ID:       projectID,
			Name:     "foobar",
			Settings: proj.Settings{},
		})
		if err != nil {
			t.Fatalf("unexpected error creating project (expected: nil, got: %v)", err)
		}

		fixtures := []sender.Request{
			{
				ID:                 ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				ProjectID:          projectID,
				SourceRequestLogID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				URL:                exampleURL,
				Method:             http.MethodPost,
				Proto:              "HTTP/1.1",
				Header: http.Header{
					"X-Foo": []string{"baz"},
				},
				Body: []byte("foo"),
				Response: &reqlog.ResponseLog{
					Proto:      "HTTP/1.1",
					Status:     "200 OK",
					StatusCode: 200,
					Header: http.Header{
						"X-Yolo": []string{"swag"},
					},
					Body: []byte("bar"),
				},
			},
			{
				ID:                 ulid.MustNew(ulid.Timestamp(time.Now())+100, ulidEntropy),
				ProjectID:          projectID,
				SourceRequestLogID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				URL:                exampleURL,
				Method:             http.MethodGet,
				Proto:              "HTTP/1.1",
				Header: http.Header{
					"X-Foo": []string{"baz"},
				},
			},
		}

		// Store fixtures.
		for _, senderReq := range fixtures {
			err = db.StoreSenderRequest(context.Background(), senderReq)
			if err != nil {
				t.Fatalf("unexpected error creating request log fixture: %v", err)
			}
		}

		filter := sender.FindRequestsFilter{
			ProjectID: projectID,
		}

		got, err := db.FindSenderRequests(context.Background(), filter, nil)
		if err != nil {
			t.Fatalf("unexpected error finding sender requests: %v", err)
		}

		// We expect the found sender requests are *reversed*, e.g. newest first.
		exp := make([]sender.Request, len(fixtures))
		for i, j := 0, len(fixtures)-1; i < j; i, j = i+1, j-1 {
			exp[i], exp[j] = fixtures[j], fixtures[i]
		}

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("sender requests not equal (-exp, +got):\n%v", diff)
		}
	})
}
