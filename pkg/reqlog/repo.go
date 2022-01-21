package reqlog

import (
	"context"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/scope"
)

type Repository interface {
	FindRequestLogs(ctx context.Context, filter FindRequestsFilter, scope *scope.Scope) ([]RequestLog, error)
	FindRequestLogByID(ctx context.Context, id ulid.ULID) (RequestLog, error)
	StoreRequestLog(ctx context.Context, reqLog RequestLog) error
	StoreResponseLog(ctx context.Context, reqLogID ulid.ULID, resLog ResponseLog) error
	ClearRequestLogs(ctx context.Context, projectID ulid.ULID) error
}
