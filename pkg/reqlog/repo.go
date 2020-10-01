package reqlog

import (
	"context"

	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/google/uuid"
)

type Repository interface {
	FindRequestLogs(ctx context.Context, opts FindRequestsOptions, scope *scope.Scope) ([]Request, error)
	FindRequestLogByID(ctx context.Context, id uuid.UUID) (Request, error)
	AddRequestLog(ctx context.Context, reqLog Request) error
	AddResponseLog(ctx context.Context, resLog Response) error
}
