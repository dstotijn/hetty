package badger

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/proj"
)

func (db *Database) UpsertProject(ctx context.Context, project proj.Project) error {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(project)
	if err != nil {
		return fmt.Errorf("badger: failed to encode project: %w", err)
	}

	err = db.badger.Update(func(txn *badger.Txn) error {
		return txn.Set(entryKey(projectPrefix, 0, project.ID[:]), buf.Bytes())
	})
	if err != nil {
		return fmt.Errorf("badger: failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) FindProjectByID(ctx context.Context, projectID ulid.ULID) (project proj.Project, err error) {
	err = db.badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get(entryKey(projectPrefix, 0, projectID[:]))
		if err != nil {
			return err
		}

		err = item.Value(func(rawProject []byte) error {
			return gob.NewDecoder(bytes.NewReader(rawProject)).Decode(&project)
		})
		if err != nil {
			return fmt.Errorf("failed to retrieve or parse project: %w", err)
		}

		return nil
	})

	if errors.Is(err, badger.ErrKeyNotFound) {
		return proj.Project{}, proj.ErrProjectNotFound
	}

	if err != nil {
		return proj.Project{}, fmt.Errorf("badger: failed to commit transaction: %w", err)
	}

	return project, nil
}

func (db *Database) DeleteProject(ctx context.Context, projectID ulid.ULID) error {
	err := db.ClearRequestLogs(ctx, projectID)
	if err != nil {
		return fmt.Errorf("badger: failed to delete project request logs: %w", err)
	}

	err = db.DeleteSenderRequests(ctx, projectID)
	if err != nil {
		return fmt.Errorf("badger: failed to delete project sender requests: %w", err)
	}

	err = db.badger.Update(func(txn *badger.Txn) error {
		return txn.Delete(entryKey(projectPrefix, 0, projectID[:]))
	})
	if err != nil {
		return fmt.Errorf("badger: failed to delete project item: %w", err)
	}

	return nil
}

func (db *Database) Projects(ctx context.Context) ([]proj.Project, error) {
	projects := make([]proj.Project, 0)

	err := db.badger.View(func(txn *badger.Txn) error {
		var rawProject []byte
		prefix := entryKey(projectPrefix, 0, nil)

		iterator := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iterator.Close()

		for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
			rawProject, err := iterator.Item().ValueCopy(rawProject)
			if err != nil {
				return fmt.Errorf("failed to copy value: %w", err)
			}

			var project proj.Project
			err = gob.NewDecoder(bytes.NewReader(rawProject)).Decode(&project)
			if err != nil {
				return fmt.Errorf("failed to decode project: %w", err)
			}

			projects = append(projects, project)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("badger: failed to commit transaction: %w", err)
	}

	return projects, nil
}
