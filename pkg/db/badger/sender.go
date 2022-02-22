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
	"github.com/dstotijn/hetty/pkg/sender"
)

func (db *Database) StoreSenderRequest(ctx context.Context, req sender.Request) error {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(req)
	if err != nil {
		return fmt.Errorf("badger: failed to encode sender request: %w", err)
	}

	entries := []*badger.Entry{
		// Sender request itself.
		{
			Key:   entryKey(senderReqPrefix, 0, req.ID[:]),
			Value: buf.Bytes(),
		},
		// Index by project ID.
		{
			Key: entryKey(senderReqPrefix, senderReqProjectIDIndex, append(req.ProjectID[:], req.ID[:]...)),
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

func (db *Database) FindSenderRequestByID(ctx context.Context, senderReqID ulid.ULID) (sender.Request, error) {
	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	req, err := getSenderRequestWithResponseLog(txn, senderReqID)
	if err != nil {
		return sender.Request{}, fmt.Errorf("badger: failed to get sender request: %w", err)
	}

	return req, nil
}

func (db *Database) FindSenderRequests(ctx context.Context, filter sender.FindRequestsFilter, scope *scope.Scope) ([]sender.Request, error) {
	if filter.ProjectID.Compare(ulid.ULID{}) == 0 {
		return nil, sender.ErrProjectIDMustBeSet
	}

	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	senderReqIDs, err := findSenderRequestIDsByProjectID(txn, filter.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("badger: failed to find sender request IDs: %w", err)
	}

	senderReqs := make([]sender.Request, 0, len(senderReqIDs))

	for _, id := range senderReqIDs {
		senderReq, err := getSenderRequestWithResponseLog(txn, id)
		if err != nil {
			return nil, fmt.Errorf("badger: failed to get sender request (id: %v): %w", id.String(), err)
		}

		if filter.OnlyInScope {
			if !senderReq.MatchScope(scope) {
				continue
			}
		}

		// Filter by search expression.
		// TODO: Once pagination is introduced, this filter logic should be done
		// as items are retrieved (e.g. when using a `badger.Iterator`).
		if filter.SearchExpr != nil {
			match, err := senderReq.Matches(filter.SearchExpr)
			if err != nil {
				return nil, fmt.Errorf(
					"badger: failed to match search expression for sender request (id: %v): %w",
					id.String(), err,
				)
			}

			if !match {
				continue
			}
		}

		senderReqs = append(senderReqs, senderReq)
	}

	return senderReqs, nil
}

func (db *Database) DeleteSenderRequests(ctx context.Context, projectID ulid.ULID) error {
	// Note: this transaction is used just for reading; we use the `badger.WriteBatch`
	// API to bulk delete items.
	txn := db.badger.NewTransaction(false)
	defer txn.Discard()

	senderReqIDs, err := findSenderRequestIDsByProjectID(txn, projectID)
	if err != nil {
		return fmt.Errorf("badger: failed to find sender request IDs: %w", err)
	}

	writeBatch := db.badger.NewWriteBatch()
	defer writeBatch.Cancel()

	for _, senderReqID := range senderReqIDs {
		// Delete sender requests.
		err := writeBatch.Delete(entryKey(senderReqPrefix, 0, senderReqID[:]))
		if err != nil {
			return fmt.Errorf("badger: failed to delete sender requests: %w", err)
		}

		// Delete related response log.
		err = writeBatch.Delete(entryKey(resLogPrefix, 0, senderReqID[:]))
		if err != nil {
			return fmt.Errorf("badger: failed to delete request log: %w", err)
		}
	}

	if err := writeBatch.Flush(); err != nil {
		return fmt.Errorf("badger: failed to commit batch write: %w", err)
	}

	err = db.badger.DropPrefix(entryKey(senderReqPrefix, senderReqProjectIDIndex, projectID[:]))
	if err != nil {
		return fmt.Errorf("badger: failed to drop sender request project ID index items: %w", err)
	}

	return nil
}

func getSenderRequestWithResponseLog(txn *badger.Txn, senderReqID ulid.ULID) (sender.Request, error) {
	item, err := txn.Get(entryKey(senderReqPrefix, 0, senderReqID[:]))

	switch {
	case errors.Is(err, badger.ErrKeyNotFound):
		return sender.Request{}, sender.ErrRequestNotFound
	case err != nil:
		return sender.Request{}, fmt.Errorf("failed to lookup sender request item: %w", err)
	}

	req := sender.Request{
		ID: senderReqID,
	}

	err = item.Value(func(rawSenderReq []byte) error {
		err = gob.NewDecoder(bytes.NewReader(rawSenderReq)).Decode(&req)
		if err != nil {
			return fmt.Errorf("failed to decode sender request: %w", err)
		}

		return nil
	})
	if err != nil {
		return sender.Request{}, fmt.Errorf("failed to retrieve or parse sender request value: %w", err)
	}

	item, err = txn.Get(entryKey(resLogPrefix, 0, senderReqID[:]))

	if errors.Is(err, badger.ErrKeyNotFound) {
		return req, nil
	}

	if err != nil {
		return sender.Request{}, fmt.Errorf("failed to get response log: %w", err)
	}

	err = item.Value(func(rawReslog []byte) error {
		var resLog reqlog.ResponseLog
		err = gob.NewDecoder(bytes.NewReader(rawReslog)).Decode(&resLog)
		if err != nil {
			return fmt.Errorf("failed to decode response log: %w", err)
		}

		req.Response = &resLog

		return nil
	})
	if err != nil {
		return sender.Request{}, fmt.Errorf("failed to retrieve or parse response log value: %w", err)
	}

	return req, nil
}

func findSenderRequestIDsByProjectID(txn *badger.Txn, projectID ulid.ULID) ([]ulid.ULID, error) {
	senderReqIDs := make([]ulid.ULID, 0)
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Reverse = true
	iterator := txn.NewIterator(opts)
	defer iterator.Close()

	var projectIndexKey []byte

	prefix := entryKey(senderReqPrefix, senderReqProjectIDIndex, projectID[:])

	for iterator.Seek(append(prefix, 255)); iterator.ValidForPrefix(prefix); iterator.Next() {
		projectIndexKey = iterator.Item().KeyCopy(projectIndexKey)

		var id ulid.ULID
		// The request log ID starts *after* the first 2 prefix and index bytes
		// and the 16 byte project ID.
		if err := id.UnmarshalBinary(projectIndexKey[18:]); err != nil {
			return nil, fmt.Errorf("failed to parse sender request ID: %w", err)
		}

		senderReqIDs = append(senderReqIDs, id)
	}

	return senderReqIDs, nil
}
