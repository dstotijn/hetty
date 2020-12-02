package reqlog

import (
	"context"
	"net/http"
	"time"

	"github.com/dstotijn/hetty/pkg/scope"
)

type RepositoryProvider interface {
	Repository() Repository
}

type Repository interface {
	FindRequestLogs(ctx context.Context, filter FindRequestsFilter, scope *scope.Scope) ([]Request, error)
	FindRequestLogByID(ctx context.Context, id int64) (Request, error)
	AddRequestLog(ctx context.Context, req http.Request, body []byte, timestamp time.Time) (*Request, error)
	AddResponseLog(ctx context.Context, reqID int64, res http.Response, body []byte, timestamp time.Time) (*Response, error)
	ClearRequestLogs(ctx context.Context) error
	UpsertSettings(ctx context.Context, module string, settings interface{}) error
	FindSettingsByModule(ctx context.Context, module string, settings interface{}) error
}
