package api

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"fmt"

	"github.com/dstotijn/gurp/pkg/reqlog"
)

type Resolver struct {
	RequestLogService *reqlog.Service
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *queryResolver) GetHTTPRequests(ctx context.Context) ([]HTTPRequest, error) {
	logs := r.RequestLogService.Requests()
	reqs := make([]HTTPRequest, len(logs))

	for i, log := range logs {
		method := HTTPMethod(log.Request.Method)
		if !method.IsValid() {
			return nil, fmt.Errorf("request has invalid method: %v", method)
		}

		reqs[i] = HTTPRequest{
			URL:       log.Request.URL.String(),
			Method:    method,
			Timestamp: log.Timestamp,
		}

		if len(log.Body) > 0 {
			reqBody := string(log.Body)
			reqs[i].Body = &reqBody
		}

		if log.Response != nil {
			reqs[i].Response = &HTTPResponse{
				StatusCode: log.Response.Response.StatusCode,
			}
			if len(log.Response.Body) > 0 {
				resBody := string(log.Response.Body)
				reqs[i].Response.Body = &resBody
			}
		}
	}

	return reqs, nil
}
