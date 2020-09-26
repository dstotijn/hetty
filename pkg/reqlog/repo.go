package reqlog

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	FindAllRequestLogs(ctx context.Context) ([]Request, error)
	FindRequestLogByID(ctx context.Context, id uuid.UUID) (Request, error)
	AddRequestLog(ctx context.Context, reqLog Request) error
	AddResponseLog(ctx context.Context, resLog Response) error
}
