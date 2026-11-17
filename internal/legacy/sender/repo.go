package sender

import (
	"context"

	"github.com/oklog/ulid"

	"github.com/pocikode/ratatoskr/internal/legacy/scope"
)

type Repository interface {
	FindSenderRequestByID(ctx context.Context, projectID, id ulid.ULID) (Request, error)
	FindSenderRequests(ctx context.Context, filter FindRequestsFilter, scope *scope.Scope) ([]Request, error)
	StoreSenderRequest(ctx context.Context, req Request) error
	DeleteSenderRequests(ctx context.Context, projectID ulid.ULID) error
}
