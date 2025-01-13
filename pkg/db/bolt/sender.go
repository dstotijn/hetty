package bolt

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"

	"github.com/oklog/ulid"
	bolt "go.etcd.io/bbolt"

	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/sender"
)

var ErrSenderRequestsBucketNotFound = errors.New("bolt: sender requests bucket not found")

var senderReqsBucketName = []byte("sender_requests")

func senderReqsBucket(tx *bolt.Tx, projectID ulid.ULID) (*bolt.Bucket, error) {
	pb, err := projectBucket(tx, projectID[:])
	if err != nil {
		return nil, err
	}

	b := pb.Bucket(senderReqsBucketName)
	if b == nil {
		return nil, ErrSenderRequestsBucketNotFound
	}

	return b, nil
}

func (db *Database) StoreSenderRequest(ctx context.Context, req sender.Request) error {
	buf := bytes.Buffer{}

	err := gob.NewEncoder(&buf).Encode(req)
	if err != nil {
		return fmt.Errorf("bolt: failed to encode sender request: %w", err)
	}

	err = db.bolt.Update(func(tx *bolt.Tx) error {
		senderReqsBucket, err := senderReqsBucket(tx, req.ProjectID)
		if err != nil {
			return fmt.Errorf("failed to get sender requests bucket: %w", err)
		}

		err = senderReqsBucket.Put(req.ID[:], buf.Bytes())
		if err != nil {
			return fmt.Errorf("failed to put sender request: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("bolt: failed to commit transaction: %w", err)
	}

	return nil
}

func (db *Database) FindSenderRequestByID(ctx context.Context, projectID, senderReqID ulid.ULID) (req sender.Request, err error) {
	if projectID.Compare(ulid.ULID{}) == 0 {
		return sender.Request{}, sender.ErrProjectIDMustBeSet
	}

	err = db.bolt.View(func(tx *bolt.Tx) error {
		senderReqsBucket, err := senderReqsBucket(tx, projectID)
		if err != nil {
			return fmt.Errorf("failed to get sender requests bucket: %w", err)
		}

		rawSenderReq := senderReqsBucket.Get(senderReqID[:])
		if rawSenderReq == nil {
			return sender.ErrRequestNotFound
		}

		err = gob.NewDecoder(bytes.NewReader(rawSenderReq)).Decode(&req)
		if err != nil {
			return fmt.Errorf("failed to decode sender request: %w", err)
		}

		return nil
	})
	if err != nil {
		return sender.Request{}, fmt.Errorf("bolt: failed to commit transaction: %w", err)
	}

	return req, nil
}

func (db *Database) FindSenderRequests(ctx context.Context, filter sender.FindRequestsFilter, scope *scope.Scope) (reqs []sender.Request, err error) {
	if filter.ProjectID.Compare(ulid.ULID{}) == 0 {
		return nil, sender.ErrProjectIDMustBeSet
	}

	tx, err := db.bolt.Begin(false)
	if err != nil {
		return nil, fmt.Errorf("bolt: failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	b, err := senderReqsBucket(tx, filter.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sender requests bucket: %w", err)
	}

	err = b.ForEach(func(senderReqID, rawSenderReq []byte) error {
		var req sender.Request
		err = gob.NewDecoder(bytes.NewReader(rawSenderReq)).Decode(&req)
		if err != nil {
			return fmt.Errorf("failed to decode sender request: %w", err)
		}

		if filter.OnlyInScope {
			if !req.MatchScope(scope) {
				return nil
			}
		}

		// Filter by search expression. TODO: Once pagination is introduced,
		// this filter logic should be done as items are retrieved.
		if filter.SearchExpr != nil {
			match, err := req.Matches(filter.SearchExpr)
			if err != nil {
				return fmt.Errorf(
					"bolt: failed to match search expression for sender request (id: %v): %w",
					senderReqID, err,
				)
			}

			if !match {
				return nil
			}
		}

		reqs = append(reqs, req)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("bolt: failed to commit transaction: %w", err)
	}

	// Reverse items, so newest requests appear first.
	for i, j := 0, len(reqs)-1; i < j; i, j = i+1, j-1 {
		reqs[i], reqs[j] = reqs[j], reqs[i]
	}

	return reqs, nil
}

func (db *Database) DeleteSenderRequests(ctx context.Context, projectID ulid.ULID) error {
	err := db.bolt.Update(func(tx *bolt.Tx) error {
		senderReqsBucket, err := senderReqsBucket(tx, projectID)
		if err != nil {
			return fmt.Errorf("failed to get sender requests bucket: %w", err)
		}

		err = senderReqsBucket.DeleteBucket(senderReqsBucketName)
		if err != nil {
			return fmt.Errorf("failed to delete sender requests bucket: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("bolt: failed to commit transaction: %w", err)
	}

	return nil
}
