package bolt_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"math/rand"
	"regexp"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/oklog/ulid"
	"go.etcd.io/bbolt"

	"github.com/dstotijn/hetty/pkg/db/bolt"
	"github.com/dstotijn/hetty/pkg/filter"
	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/scope"
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

	path := t.TempDir() + "bolt.db"
	boltDB, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		t.Fatalf("failed to open bolt database: %v", err)
	}

	db, err := bolt.DatabaseFromBoltDB(boltDB)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	searchExpr, err := filter.ParseQuery("foo AND bar OR NOT baz")
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

	err = db.UpsertProject(context.Background(), exp)
	if err != nil {
		t.Fatalf("unexpected error storing project: %v", err)
	}

	var rawProject []byte

	err = boltDB.View(func(tx *bbolt.Tx) error {
		rawProject = tx.Bucket([]byte("projects")).Bucket(exp.ID[:]).Get([]byte("project"))
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error retrieving project from database: %v", err)
	}
	if rawProject == nil {
		t.Fatalf("expected raw project to be retrieved, got: nil")
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

		path := t.TempDir() + "bolt.db"
		boltDB, err := bbolt.Open(path, 0o600, nil)
		if err != nil {
			t.Fatalf("failed to open bolt database: %v", err)
		}

		db, err := bolt.DatabaseFromBoltDB(boltDB)
		if err != nil {
			t.Fatalf("failed to create database: %v", err)
		}
		defer db.Close()

		exp := proj.Project{
			ID: ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
		}

		buf := bytes.Buffer{}

		err = gob.NewEncoder(&buf).Encode(exp)
		if err != nil {
			t.Fatalf("unexpected error encoding project: %v", err)
		}

		err = boltDB.Update(func(tx *bbolt.Tx) error {
			b, err := tx.Bucket([]byte("projects")).CreateBucket(exp.ID[:])
			if err != nil {
				return err
			}
			return b.Put([]byte("project"), buf.Bytes())
		})
		if err != nil {
			t.Fatalf("unexpected error setting project: %v", err)
		}

		got, err := db.FindProjectByID(context.Background(), exp.ID)
		if err != nil {
			t.Fatalf("unexpected error finding project: %v", err)
		}

		if diff := cmp.Diff(exp, got, cmpopts.IgnoreUnexported(proj.Project{})); diff != "" {
			t.Fatalf("project not equal (-exp, +got):\n%v", diff)
		}
	})

	t.Run("project not found", func(t *testing.T) {
		t.Parallel()

		path := t.TempDir() + "bolt.db"
		boltDB, err := bbolt.Open(path, 0o600, nil)
		if err != nil {
			t.Fatalf("failed to open bolt database: %v", err)
		}

		db, err := bolt.DatabaseFromBoltDB(boltDB)
		if err != nil {
			t.Fatalf("failed to create database: %v", err)
		}
		defer db.Close()

		projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)

		_, err = db.FindProjectByID(context.Background(), projectID)
		if !errors.Is(err, proj.ErrProjectNotFound) {
			t.Fatalf("expected `proj.ErrProjectNotFound`, got: %v", err)
		}
	})
}

func TestDeleteProject(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "bolt.db"
	boltDB, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		t.Fatalf("failed to open bolt database: %v", err)
	}

	db, err := bolt.DatabaseFromBoltDB(boltDB)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	// Insert test fixture.
	projectID := ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	err = db.UpsertProject(context.Background(), proj.Project{
		ID: projectID,
	})
	if err != nil {
		t.Fatalf("unexpected error storing project: %v", err)
	}

	err = db.DeleteProject(context.Background(), projectID)
	if err != nil {
		t.Fatalf("unexpected error deleting project: %v", err)
	}

	var got *bbolt.Bucket
	err = boltDB.View(func(tx *bbolt.Tx) error {
		got = tx.Bucket([]byte("projects")).Bucket(projectID[:])
		return nil
	})
	if got != nil {
		t.Fatalf("expected bucket to be nil, got: %v", got)
	}
}

func TestProjects(t *testing.T) {
	t.Parallel()

	path := t.TempDir() + "bolt.db"
	boltDB, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		t.Fatalf("failed to open bolt database: %v", err)
	}

	db, err := bolt.DatabaseFromBoltDB(boltDB)
	if err != nil {
		t.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

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
		err = db.UpsertProject(context.Background(), project)
		if err != nil {
			t.Fatalf("unexpected error creating project fixture: %v", err)
		}
	}

	got, err := db.Projects(context.Background())
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
