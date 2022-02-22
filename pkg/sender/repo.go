package sender

import (
	"context"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
)

type Repository interface {
	FindSenderRequestByID(ctx context.Context, id ulid.ULID) (Request, error)
	FindSenderRequests(ctx context.Context, filter FindRequestsFilter, scope *scope.Scope) ([]Request, error)
	StoreSenderRequest(ctx context.Context, req Request) error
	StoreResponseLog(ctx context.Context, reqLogID ulid.ULID, resLog reqlog.ResponseLog) error
	DeleteSenderRequests(ctx context.Context, projectID ulid.ULID) error
}
