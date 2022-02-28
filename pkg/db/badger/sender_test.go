package badger_test

import (
	"context"
	"errors"
	"math/rand"
	"net/http"
	"net/url"
	"testing"
	"time"

	badgerdb "github.com/dgraph-io/badger/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/db/badger"
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

func TestFindRequestByID(t *testing.T) {
	t.Parallel()

	database, err := badger.OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open badger database: %v", err)
	}
	defer database.Close()

	// See: https://go.dev/blog/subtests#cleaning-up-after-a-group-of-parallel-tests
	t.Run("group", func(t *testing.T) {
		t.Run("sender request not found", func(t *testing.T) {
			t.Parallel()

			_, err := database.FindSenderRequestByID(context.Background(), ulid.ULID{})
			if !errors.Is(err, sender.ErrRequestNotFound) {
				t.Fatalf("expected `sender.ErrRequestNotFound`, got: %v", err)
			}
		})

		t.Run("sender request found", func(t *testing.T) {
			t.Parallel()

			exp := sender.Request{
				ID:                 ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				ProjectID:          ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				SourceRequestLogID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),

				URL:    exampleURL,
				Method: http.MethodGet,
				Proto:  sender.HTTPProto2,
				Header: http.Header{
					"X-Foo": []string{"bar"},
				},
				Body: []byte("foo"),
			}

			err := database.StoreSenderRequest(context.Background(), exp)
			if err != nil {
				t.Fatalf("unexpected error (expected: nil, got: %v)", err)
			}

			resLog := reqlog.ResponseLog{
				Proto:      "HTTP/2.0",
				Status:     "200 OK",
				StatusCode: 200,
				Header: http.Header{
					"X-Yolo": []string{"swag"},
				},
				Body: []byte("bar"),
			}

			err = database.StoreResponseLog(context.Background(), exp.ID, resLog)
			if err != nil {
				t.Fatalf("unexpected error (expected: nil, got: %v)", err)
			}

			exp.Response = &resLog

			got, err := database.FindSenderRequestByID(context.Background(), exp.ID)
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

		database, err := badger.OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}
		defer database.Close()

		filter := sender.FindRequestsFilter{}

		_, err = database.FindSenderRequests(context.Background(), filter, nil)
		if !errors.Is(err, sender.ErrProjectIDMustBeSet) {
			t.Fatalf("expected `sender.ErrProjectIDMustBeSet`, got: %v", err)
		}
	})

	t.Run("returns sender requests and related response logs", func(t *testing.T) {
		t.Parallel()

		database, err := badger.OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}
		defer database.Close()

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

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
			err = database.StoreSenderRequest(context.Background(), senderReq)
			if err != nil {
				t.Fatalf("unexpected error creating request log fixture: %v", err)
			}

			if senderReq.Response != nil {
				err = database.StoreResponseLog(context.Background(), senderReq.ID, *senderReq.Response)
				if err != nil {
					t.Fatalf("unexpected error creating response log fixture: %v", err)
				}
			}
		}

		filter := sender.FindRequestsFilter{
			ProjectID: projectID,
		}

		got, err := database.FindSenderRequests(context.Background(), filter, nil)
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
