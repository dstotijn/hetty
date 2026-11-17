package reqlog

import (
	"context"

	"github.com/oklog/ulid"

	"github.com/pocikode/ratatoskr/internal/legacy/scope"
)

type Repository interface {
	FindRequestLogs(ctx context.Context, filter FindRequestsFilter, scope *scope.Scope) ([]RequestLog, error)
	FindRequestLogByID(ctx context.Context, projectID, id ulid.ULID) (RequestLog, error)
	StoreRequestLog(ctx context.Context, reqLog RequestLog) error
	StoreResponseLog(ctx context.Context, projectID, reqLogID ulid.ULID, resLog ResponseLog) error
	ClearRequestLogs(ctx context.Context, projectID ulid.ULID) error
}
