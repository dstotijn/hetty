package reqlog

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/scope"
)

type contextKey int

const LogBypassedKey contextKey = 0

var ErrRequestNotFound = errors.New("reqlog: request not found")

type Request struct {
	ID        int64
	Request   http.Request
	Body      []byte
	Timestamp time.Time
	Response  *Response
}

type Response struct {
	ID        int64
	RequestID int64
	Response  http.Response
	Body      []byte
	Timestamp time.Time
}

type Service struct {
	BypassOutOfScopeRequests bool

	scope *scope.Scope
	repo  Repository
}

type FindRequestsOptions struct {
	OmitOutOfScope bool
}

type Config struct {
	Scope                    *scope.Scope
	Repository               Repository
	BypassOutOfScopeRequests bool
}

func NewService(cfg Config) *Service {
	return &Service{
		scope:                    cfg.Scope,
		repo:                     cfg.Repository,
		BypassOutOfScopeRequests: cfg.BypassOutOfScopeRequests,
	}
}

func (svc *Service) FindRequests(ctx context.Context, opts FindRequestsOptions) ([]Request, error) {
	var scope *scope.Scope
	if opts.OmitOutOfScope {
		scope = svc.scope
	}

	return svc.repo.FindRequestLogs(ctx, opts, scope)
}

func (svc *Service) FindRequestLogByID(ctx context.Context, id int64) (Request, error) {
	return svc.repo.FindRequestLogByID(ctx, id)
}

func (svc *Service) addRequest(
	ctx context.Context,
	req http.Request,
	body []byte,
	timestamp time.Time,
) (*Request, error) {
	return svc.repo.AddRequestLog(ctx, req, body, timestamp)
}

func (svc *Service) addResponse(
	ctx context.Context,
	reqID int64,
	res http.Response,
	body []byte,
	timestamp time.Time,
) (*Response, error) {
	if res.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("reqlog: could not create gzip reader: %v", err)
		}
		defer gzipReader.Close()
		body, err = ioutil.ReadAll(gzipReader)
		if err != nil {
			return nil, fmt.Errorf("reqlog: could not read gzipped response body: %v", err)
		}
	}

	return svc.repo.AddResponseLog(ctx, reqID, res, body, timestamp)
}

func (svc *Service) RequestModifier(next proxy.RequestModifyFunc) proxy.RequestModifyFunc {
	return func(req *http.Request) {
		now := time.Now()
		next(req)

		clone := req.Clone(req.Context())
		var body []byte
		if req.Body != nil {
			// TODO: Use io.LimitReader.
			var err error
			body, err = ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("[ERROR] Could not read request body for logging: %v", err)
				return
			}
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		// Bypass logging if this setting is enabled and the incoming request
		// doens't match any rules of the scope.
		if svc.BypassOutOfScopeRequests && !svc.scope.Match(clone, body) {
			ctx := context.WithValue(req.Context(), LogBypassedKey, true)
			*req = *req.WithContext(ctx)
			return
		}

		reqLog, err := svc.addRequest(req.Context(), *clone, body, now)
		if err != nil {
			log.Printf("[ERROR] Could not store request log: %v", err)
			return
		}
		ctx := context.WithValue(req.Context(), proxy.ReqIDKey, reqLog.ID)
		*req = *req.WithContext(ctx)
	}
}

func (svc *Service) ResponseModifier(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
	return func(res *http.Response) error {
		now := time.Now()
		if err := next(res); err != nil {
			return err
		}

		if bypassed, _ := res.Request.Context().Value(LogBypassedKey).(bool); bypassed {
			return nil
		}

		reqID, _ := res.Request.Context().Value(proxy.ReqIDKey).(int64)
		if reqID == 0 {
			return errors.New("reqlog: request is missing ID")
		}

		clone := *res

		// TODO: Use io.LimitReader.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reqlog: could not read response body: %v", err)
		}
		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		go func() {
			if _, err := svc.addResponse(context.Background(), reqID, clone, body, now); err != nil {
				log.Printf("[ERROR] Could not store response log: %v", err)
			}
		}()

		return nil
	}
}
