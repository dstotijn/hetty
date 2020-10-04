package api

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"fmt"

	"github.com/dstotijn/hetty/pkg/reqlog"
)

type Resolver struct {
	RequestLogService *reqlog.Service
}

type queryResolver struct{ *Resolver }

func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

func (r *queryResolver) HTTPRequestLogs(ctx context.Context) ([]HTTPRequestLog, error) {
	opts := reqlog.FindRequestsOptions{OmitOutOfScope: false}
	reqs, err := r.RequestLogService.FindRequests(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("could not query repository for requests: %v", err)
	}
	logs := make([]HTTPRequestLog, len(reqs))

	for i, req := range reqs {
		req, err := parseRequestLog(req)
		if err != nil {
			return nil, err
		}
		logs[i] = req
	}

	return logs, nil
}

func (r *queryResolver) HTTPRequestLog(ctx context.Context, id int64) (*HTTPRequestLog, error) {
	log, err := r.RequestLogService.FindRequestLogByID(ctx, id)
	if err == reqlog.ErrRequestNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get request by ID: %v", err)
	}
	req, err := parseRequestLog(log)
	if err != nil {
		return nil, err
	}

	return &req, nil
}

func parseRequestLog(req reqlog.Request) (HTTPRequestLog, error) {
	method := HTTPMethod(req.Request.Method)
	if !method.IsValid() {
		return HTTPRequestLog{}, fmt.Errorf("request has invalid method: %v", method)
	}

	log := HTTPRequestLog{
		ID:        req.ID,
		URL:       req.Request.URL.String(),
		Proto:     req.Request.Proto,
		Method:    method,
		Timestamp: req.Timestamp,
	}

	if len(req.Body) > 0 {
		reqBody := string(req.Body)
		log.Body = &reqBody
	}

	if req.Request.Header != nil {
		log.Headers = make([]HTTPHeader, 0)
		for key, values := range req.Request.Header {
			for _, value := range values {
				log.Headers = append(log.Headers, HTTPHeader{
					Key:   key,
					Value: value,
				})
			}
		}
	}

	if req.Response != nil {
		log.Response = &HTTPResponseLog{
			RequestID:  req.ID,
			Proto:      req.Response.Response.Proto,
			Status:     req.Response.Response.Status,
			StatusCode: req.Response.Response.StatusCode,
		}
		if len(req.Response.Body) > 0 {
			resBody := string(req.Response.Body)
			log.Response.Body = &resBody
		}
		if req.Response.Response.Header != nil {
			log.Response.Headers = make([]HTTPHeader, 0)
			for key, values := range req.Response.Response.Header {
				for _, value := range values {
					log.Response.Headers = append(log.Response.Headers, HTTPHeader{
						Key:   key,
						Value: value,
					})
				}
			}
		}
	}

	return log, nil
}
