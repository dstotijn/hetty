package api

import (
	"context"
	"fmt"

	"github.com/dstotijn/gurp/pkg/reqlog"
)

type Resolver struct {
	RequestLogStore *reqlog.RequestLogStore
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *queryResolver) GetRequests(ctx context.Context) ([]Request, error) {
	reqs := r.RequestLogStore.Requests()
	resp := make([]Request, len(reqs))

	for i := range resp {
		method := HTTPMethod(reqs[i].Request.Method)
		if !method.IsValid() {
			return nil, fmt.Errorf("request has invalid method: %v", method)
		}
		resp[i] = Request{
			URL:    reqs[i].Request.URL.String(),
			Method: method,
		}
	}

	return resp, nil
}
