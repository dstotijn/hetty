package reqlog

import (
	"context"
	"net/http"
	"time"

	"github.com/dstotijn/hetty/pkg/scope"
)

type Repository interface {
	FindRequestLogs(ctx context.Context, opts FindRequestsOptions, scope *scope.Scope) ([]Request, error)
	FindRequestLogByID(ctx context.Context, id int64) (Request, error)
	AddRequestLog(ctx context.Context, req http.Request, body []byte, timestamp time.Time) (*Request, error)
	AddResponseLog(ctx context.Context, reqID int64, res http.Response, body []byte, timestamp time.Time) (*Response, error)
}
