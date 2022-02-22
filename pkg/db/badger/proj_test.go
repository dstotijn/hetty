package badger

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"math/rand"
	"regexp"
	"testing"
	"time"

	badgerdb "github.com/dgraph-io/badger/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/search"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

var regexpCompareOpt = cmp.Comparer(func(x, y *regexp.Regexp) bool {
	switch {
	case x == nil && y == nil:
		return true
	case x == nil || y == nil:
		return false
	default:
		return x.String() == y.String()
	}
})

func TestUpsertProject(t *testing.T) {
	t.Parallel()

	badgerDB, err := badgerdb.Open(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open badger database: %v", err)
	}

	database := DatabaseFromBadgerDB(badgerDB)
	defer database.Close()

	searchExpr, err := search.ParseQuery("foo AND bar OR NOT baz")
	if err != nil {
		t.Fatalf("unexpected error (expected: nil, got: %v)", err)
	}

	exp := proj.Project{
		ID:   ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
		Name: "foobar",
		Settings: proj.Settings{
			ReqLogBypassOutOfScope: true,
			ReqLogOnlyFindInScope:  true,
			ReqLogSearchExpr:       searchExpr,
			ScopeRules: []scope.Rule{
				{
					URL: regexp.MustCompile("^https://(.*)example.com(.*)$"),
					Header: scope.Header{
						Key:   regexp.MustCompile("^X-Foo(.*)$"),
						Value: regexp.MustCompile("^foo(.*)$"),
					},
					Body: regexp.MustCompile("^foo(.*)"),
				},
			},
		},
	}

	err = database.UpsertProject(context.Background(), exp)
	if err != nil {
		t.Fatalf("unexpected error storing project: %v", err)
	}

	var rawProject []byte

	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		item, err := txn.Get(entryKey(projectPrefix, 0, exp.ID[:]))
		if err != nil {
			return err
		}

		rawProject, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		t.Fatalf("unexpected error retrieving project from database: %v", err)
	}

	got := proj.Project{}

	err = gob.NewDecoder(bytes.NewReader(rawProject)).Decode(&got)
	if err != nil {
		t.Fatalf("unexpected error decoding project: %v", err)
	}

	if diff := cmp.Diff(exp, got, regexpCompareOpt, cmpopts.IgnoreUnexported(proj.Project{})); diff != "" {
		t.Fatalf("project not equal (-exp, +got):\n%v", diff)
	}
}

func TestFindProjectByID(t *testing.T) {
	t.Parallel()

	t.Run("existing project", func(t *testing.T) {
		t.Parallel()

		badgerDB, err := badgerdb.Open(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}

		database := DatabaseFromBadgerDB(badgerDB)
		defer database.Close()

		exp := proj.Project{
			ID:       ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
			Name:     "foobar",
			Settings: proj.Settings{},
		}

		buf := bytes.Buffer{}

		err = gob.NewEncoder(&buf).Encode(exp)
		if err != nil {
			t.Fatalf("unexpected error encoding project: %v", err)
		}

		err = badgerDB.Update(func(txn *badgerdb.Txn) error {
			return txn.Set(entryKey(projectPrefix, 0, exp.ID[:]), buf.Bytes())
		})
		if err != nil {
			t.Fatalf("unexpected error setting project: %v", err)
		}

		got, err := database.FindProjectByID(context.Background(), exp.ID)
		if err != nil {
			t.Fatalf("unexpected error finding project: %v", err)
		}

		if diff := cmp.Diff(exp, got, cmpopts.IgnoreUnexported(proj.Project{})); diff != "" {
			t.Fatalf("project not equal (-exp, +got):\n%v", diff)
		}
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		database, err := OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
		if err != nil {
			t.Fatalf("failed to open badger database: %v", err)
		}
		defer database.Close()

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

		_, err = database.FindProjectByID(context.Background(), projectID)
		if !errors.Is(err, proj.ErrProjectNotFound) {
			t.Fatalf("expected `proj.ErrProjectNotFound`, got: %v", err)
		}
	})
}

func TestDeleteProject(t *testing.T) {
	t.Parallel()

	badgerDB, err := badgerdb.Open(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open badger database: %v", err)
	}

	database := DatabaseFromBadgerDB(badgerDB)
	defer database.Close()

	// Store fixtures.
	projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	reqLogID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	senderReqID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

	err = badgerDB.Update(func(txn *badgerdb.Txn) error {
		// Project item.
		if err := txn.Set(entryKey(projectPrefix, 0, projectID[:]), nil); err != nil {
			return err
		}

		// Sender request items.
		if err := txn.Set(entryKey(senderReqPrefix, 0, senderReqID[:]), nil); err != nil {
			return err
		}
		if err := txn.Set(entryKey(resLogPrefix, 0, senderReqID[:]), nil); err != nil {
			return err
		}
		err := txn.Set(entryKey(senderReqPrefix, senderReqProjectIDIndex, append(projectID[:], senderReqID[:]...)), nil)
		if err != nil {
			return err
		}

		// Request log items.
		if err := txn.Set(entryKey(reqLogPrefix, 0, reqLogID[:]), nil); err != nil {
			return err
		}
		if err := txn.Set(entryKey(resLogPrefix, 0, reqLogID[:]), nil); err != nil {
			return err
		}
		err = txn.Set(entryKey(reqLogPrefix, reqLogProjectIDIndex, append(projectID[:], reqLogID[:]...)), nil)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error creating fixtures: %v", err)
	}

	err = database.DeleteProject(context.Background(), projectID)
	if err != nil {
		t.Fatalf("unexpected error deleting project: %v", err)
	}

	// Assert project key was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(projectPrefix, 0, projectID[:]))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert request log item was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(reqLogPrefix, 0, reqLogID[:]))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert response log item related to request log was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(resLogPrefix, 0, reqLogID[:]))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert request log project ID index key was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(reqLogPrefix, reqLogProjectIDIndex, append(projectID[:], reqLogID[:]...)))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert sender request item was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(senderReqPrefix, 0, senderReqID[:]))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert response log item related to sender request was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(resLogPrefix, 0, senderReqID[:]))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}

	// Assert sender request project ID index key was deleted.
	err = badgerDB.View(func(txn *badgerdb.Txn) error {
		_, err := txn.Get(entryKey(senderReqPrefix, senderReqProjectIDIndex, append(projectID[:], senderReqID[:]...)))
		return err
	})
	if !errors.Is(err, badgerdb.ErrKeyNotFound) {
		t.Fatalf("expected `badger.ErrKeyNotFound`, got: %v", err)
	}
}

func TestProjects(t *testing.T) {
	t.Parallel()

	database, err := OpenDatabase(badgerdb.DefaultOptions("").WithInMemory(true))
	if err != nil {
		t.Fatalf("failed to open badger database: %v", err)
	}
	defer database.Close()

	exp := []proj.Project{
		{
			ID:   ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
			Name: "one",
		},
		{
			ID:   ulid.MustNew(ulid.Timestamp(time.Now())+100, ulidEntropy),
			Name: "two",
		},
	}

	// Store fixtures.
	for _, project := range exp {
		err = database.UpsertProject(context.Background(), project)
		if err != nil {
			t.Fatalf("unexpected error creating project fixture: %v", err)
		}
	}

	got, err := database.Projects(context.Background())
	if err != nil {
		t.Fatalf("unexpected error finding projects: %v", err)
	}

	if len(exp) != len(got) {
		t.Fatalf("expected %v projects, got: %v", len(exp), len(got))
	}

	if diff := cmp.Diff(exp, got, cmpopts.IgnoreUnexported(proj.Project{})); diff != "" {
		t.Fatalf("projects not equal (-exp, +got):\n%v", diff)
	}
}
