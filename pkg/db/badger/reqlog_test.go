package badger

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	badgerdb "github.com/dgraph-io/badger/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/reqlog"
)

func TestFindRequestLogs(t *testing.T) {
	t.Parallel()

	t.Run("without project ID in filter", func(t *testing.T) {
		t.Parallel()

		database, err := OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}
		defer database.Close()

		filter := reqlog.FindRequestsFilter{}

		_, err = database.FindRequestLogs(context.Background(), filter, nil)
		if !errors.Is(err, reqlog.ErrProjectIDMustBeSet) {
			t.Fatalf("expected `reqlog.ErrProjectIDMustBeSet`, got: %v", err)
		}
	})

	t.Run("returns request logs and related response logs", func(t *testing.T) {
		t.Parallel()

		database, err := OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}
		defer database.Close()

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

		exp := []reqlog.RequestLog{
			{
				ID:        ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
				ProjectID: projectID,
				URL:       mustParseURL(t, "https://example.com/foobar"),
				Method:    http.MethodPost,
				Proto:     "HTTP/1.1",
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
				ID:        ulid.MustNew(ulid.Timestamp(time.Now())+100, ulidEntropy),
				ProjectID: projectID,
				URL:       mustParseURL(t, "https://example.com/foo?bar=baz"),
				Method:    http.MethodGet,
				Proto:     "HTTP/1.1",
				Header: http.Header{
					"X-Foo": []string{"baz"},
				},
			},
		}

		// Store fixtures.
		for _, reqLog := range exp {
			err = database.StoreRequestLog(context.Background(), reqLog)
			if err != nil {
				t.Fatalf("unexpected error creating request log fixture: %v", err)
			}

			if reqLog.Response != nil {
				err = database.StoreResponseLog(context.Background(), reqLog.ID, *reqLog.Response)
				if err != nil {
					t.Fatalf("unexpected error creating response log fixture: %v", err)
				}
			}
		}

		filter := reqlog.FindRequestsFilter{
			ProjectID: projectID,
		}

		got, err := database.FindRequestLogs(context.Background(), filter, nil)
		if err != nil {
			t.Fatalf("unexpected error finding request logs: %v", err)
		}

		if diff := cmp.Diff(exp, got); diff != "" {
			t.Fatalf("request logs not equal (-exp, +got):\n%v", diff)
		}
	})
}

func mustParseURL(t *testing.T, s string) *url.URL {
	t.Helper()

	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u
}
