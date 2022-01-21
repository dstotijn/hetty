package badger

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
)

func (db *Database) FindRequestLogs(ctx context.Context, filter reqlog.FindRequestsFilter, scope *scope.Scope) ([]reqlog.RequestLog, error) {
	if filter.ProjectID.Compare(ulid.ULID{}) == 0 {
		return nil, reqlog.ErrProjectIDMustBeSet
	}

	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	reqLogIDs, err := findRequestLogIDsByProjectID(txn, filter.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("badger: failed to find request log IDs: %w", err)
	}

	reqLogs := make([]reqlog.RequestLog, 0, len(reqLogIDs))

	for _, reqLogID := range reqLogIDs {
		reqLog, err := getRequestLogWithResponse(txn, reqLogID)
		if err != nil {
			return nil, fmt.Errorf("badger: failed to get request log (id: %v): %w", reqLogID.String(), err)
		}

		if filter.OnlyInScope {
			if !reqLog.MatchScope(scope) {
				continue
			}
		}

		// Filter by search expression.
		// TODO: Once pagination is introduced, this filter logic should be done
		// as items are retrieved (e.g. when using a `badger.Iterator`).
		if filter.SearchExpr != nil {
			match, err := reqLog.Matches(filter.SearchExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"badger: failed to match search expression for request log (id: %v): %w",
					reqLogID.String(), err,
				)
			}

			if !match {
				continue
			}
		}

		reqLogs = append(reqLogs, reqLog)
	}

	return reqLogs, nil
}

func getRequestLogWithResponse(txn *badger.Txn, reqLogID ulid.ULID) (reqlog.RequestLog, error) {
	item, err := txn.Get(entryKey(reqLogPrefix, 0, reqLogID[:]))
	if err != nil {
		return reqlog.RequestLog{}, fmt.Errorf("failed to lookup request log item: %w", err)
	}

	reqLog := reqlog.RequestLog{
		ID: reqLogID,
	}

	err = item.Value(func(rawReqLog []byte) error {
		err = gob.NewDecoder(bytes.NewReader(rawReqLog)).Decode(&reqLog)
		if err != nil {
			return fmt.Errorf("failed to decode request log: %w", err)
		}

		return nil
	})
	if err != nil {
		return reqlog.RequestLog{}, fmt.Errorf("failed to retrieve or parse request log value: %w", err)
	}

	item, err = txn.Get(entryKey(resLogPrefix, 0, reqLogID[:]))

	if errors.Is(err, badger.ErrKeyNotFound) {
		return reqLog, nil
	}

	if err != nil {
		return reqlog.RequestLog{}, fmt.Errorf("failed to get response log: %w", err)
	}

	err = item.Value(func(rawReslog []byte) error {
		var resLog reqlog.ResponseLog
		err = gob.NewDecoder(bytes.NewReader(rawReslog)).Decode(&resLog)
		if err != nil {
			return fmt.Errorf("failed to decode response log: %w", err)
		}

		reqLog.Response = &resLog

		return nil
	})
	if err != nil {
		return reqlog.RequestLog{}, fmt.Errorf("failed to retrieve or parse response log value: %w", err)
	}

	return reqLog, nil
}

func (db *Database) FindRequestLogByID(ctx context.Context, reqLogID ulid.ULID) (reqLog reqlog.RequestLog, err error) {
	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	reqLog, err = getRequestLogWithResponse(txn, reqLogID)
	if err != nil {
		return reqlog.RequestLog{}, fmt.Errorf("badger: failed to get request log: %w", err)
	}

	return reqLog, nil
}

func (db *Database) StoreRequestLog(ctx context.Context, reqLog reqlog.RequestLog) error {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(reqLog)
	if err != nil {
		return fmt.Errorf("badger: failed to encode request log: %w", err)
	}

	entries := []*badger.Entry{
		// Request log itself.
		{
			Key:   entryKey(reqLogPrefix, 0, reqLog.ID[:]),
			Value: buf.Bytes(),
		},
		// Index by project ID.
		{
			Key: entryKey(reqLogPrefix, reqLogProjectIDIndex, append(reqLog.ProjectID[:], reqLog.ID[:]...)),
		},
	}

	err = db.badger.Update(func(txn *badger.Txn) error {
		for i := range entries {
			err := txn.SetEntry(entries[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("badger: failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) StoreResponseLog(ctx context.Context, reqLogID ulid.ULID, resLog reqlog.ResponseLog) error {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(resLog)
	if err != nil {
		return fmt.Errorf("badger: failed to encode response log: %w", err)
	}

	err = db.badger.Update(func(txn *badger.Txn) error {
		return txn.SetEntry(&badger.Entry{
			Key:   entryKey(resLogPrefix, 0, reqLogID[:]),
			Value: buf.Bytes(),
		})
	})
	if err != nil {
		return fmt.Errorf("badger: failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) ClearRequestLogs(ctx context.Context, projectID ulid.ULID) error {
	// Note: this transaction is used just for reading; we use the `badger.WriteBatch`
	// API to bulk delete items.
	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	reqLogIDs, err := findRequestLogIDsByProjectID(txn, projectID)
	if err != nil {
		return fmt.Errorf("badger: failed to find request log IDs: %w", err)
	}

	writeBatch := db.badger.NewWriteBatch()
	defer writeBatch.Cancel()

	for _, reqLogID := range reqLogIDs {
		// Delete request logs.
		err := writeBatch.Delete(entryKey(reqLogPrefix, 0, reqLogID[:]))
		if err != nil {
			return fmt.Errorf("badger: failed to delete request log: %w", err)
		}

		// Delete related response log.
		err = writeBatch.Delete(entryKey(resLogPrefix, 0, reqLogID[:]))
		if err != nil {
			return fmt.Errorf("badger: failed to delete request log: %w", err)
		}
	}

	if err := writeBatch.Flush(); err != nil {
		return fmt.Errorf("badger: failed to commit batch write: %w", err)
	}

	err = db.badger.DropPrefix(entryKey(reqLogPrefix, reqLogProjectIDIndex, projectID[:]))
	if err != nil {
		return fmt.Errorf("badger: failed to drop request log project ID index items: %w", err)
	}

	return nil
}

func findRequestLogIDsByProjectID(txn *badger.Txn, projectID ulid.ULID) ([]ulid.ULID, error) {
	reqLogIDs := make([]ulid.ULID, 0)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	iterator := txn.NewIterator(opts)
	defer iterator.Close()

	var projectIndexKey []byte

	prefix := entryKey(reqLogPrefix, reqLogProjectIDIndex, projectID[:])

	for iterator.Seek(prefix); iterator.ValidForPrefix(prefix); iterator.Next() {
		projectIndexKey = iterator.Item().KeyCopy(projectIndexKey)

		var id ulid.ULID
		// The request log ID starts *after* the first 2 prefix and index bytes
		// and the 16 byte project ID.
		if err := id.UnmarshalBinary(projectIndexKey[18:]); err != nil {
			return nil, fmt.Errorf("failed to parse request log ID: %w", err)
		}

		reqLogIDs = append(reqLogIDs, id)
	}

	return reqLogIDs, nil
}
