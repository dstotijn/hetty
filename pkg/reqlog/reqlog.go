package reqlog

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/scope"
	"github.com/dstotijn/hetty/pkg/search"
)

type contextKey int

const LogBypassedKey contextKey = 0

var (
	ErrRequestNotFound    = errors.New("reqlog: request not found")
	ErrProjectIDMustBeSet = errors.New("reqlog: project ID must be set")
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

type RequestLog struct {
	ID        ulid.ULID
	ProjectID ulid.ULID

	URL    *url.URL
	Method string
	Proto  string
	Header http.Header
	Body   []byte

	Response *ResponseLog
}

type ResponseLog struct {
	Proto      string
	StatusCode int
	Status     string
	Header     http.Header
	Body       []byte
}

type Service struct {
	BypassOutOfScopeRequests bool
	FindReqsFilter           FindRequestsFilter
	ActiveProjectID          ulid.ULID

	scope *scope.Scope
	repo  Repository
}

type FindRequestsFilter struct {
	ProjectID   ulid.ULID
	OnlyInScope bool
	SearchExpr  search.Expression
}

type Config struct {
	Scope      *scope.Scope
	Repository Repository
}

func NewService(cfg Config) *Service {
	return &Service{
		repo:  cfg.Repository,
		scope: cfg.Scope,
	}
}

func (svc *Service) FindRequests(ctx context.Context) ([]RequestLog, error) {
	return svc.repo.FindRequestLogs(ctx, svc.FindReqsFilter, svc.scope)
}

func (svc *Service) FindRequestLogByID(ctx context.Context, id ulid.ULID) (RequestLog, error) {
	return svc.repo.FindRequestLogByID(ctx, id)
}

func (svc *Service) ClearRequests(ctx context.Context, projectID ulid.ULID) error {
	return svc.repo.ClearRequestLogs(ctx, projectID)
}

func (svc *Service) storeResponse(ctx context.Context, reqLogID ulid.ULID, res *http.Response) error {
	if res.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(res.Body)
		if err != nil {
			return fmt.Errorf("could not create gzip reader: %w", err)
		}
		defer gzipReader.Close()

		buf := &bytes.Buffer{}

		if _, err := io.Copy(buf, gzipReader); err != nil {
			return fmt.Errorf("could not read gzipped response body: %w", err)
		}

		res.Body = io.NopCloser(buf)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("could not read body: %w", err)
	}

	resLog := ResponseLog{
		Proto:      res.Proto,
		StatusCode: res.StatusCode,
		Status:     res.Status,
		Header:     res.Header,
		Body:       body,
	}

	return svc.repo.StoreResponseLog(ctx, reqLogID, resLog)
}

func (svc *Service) RequestModifier(next proxy.RequestModifyFunc) proxy.RequestModifyFunc {
	return func(req *http.Request) {
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
			clone.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		// Bypass logging if no project is active.
		if svc.ActiveProjectID.Compare(ulid.ULID{}) == 0 {
			ctx := context.WithValue(req.Context(), LogBypassedKey, true)
			*req = *req.WithContext(ctx)

			return
		}

		// Bypass logging if this setting is enabled and the incoming request
		// doesn't match any scope rules.
		if svc.BypassOutOfScopeRequests && !svc.scope.Match(clone, body) {
			ctx := context.WithValue(req.Context(), LogBypassedKey, true)
			*req = *req.WithContext(ctx)

			return
		}

		reqLog := RequestLog{
			ID:        ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
			ProjectID: svc.ActiveProjectID,
			Method:    clone.Method,
			URL:       clone.URL,
			Proto:     clone.Proto,
			Header:    clone.Header,
			Body:      body,
		}

		err := svc.repo.StoreRequestLog(req.Context(), reqLog)
		if err != nil {
			log.Printf("[ERROR] Could not store request log: %v", err)
			return
		}

		ctx := context.WithValue(req.Context(), proxy.ReqLogIDKey, reqLog.ID)
		*req = *req.WithContext(ctx)
	}
}

func (svc *Service) ResponseModifier(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		if bypassed, _ := res.Request.Context().Value(LogBypassedKey).(bool); bypassed {
			return nil
		}

		reqLogID, ok := res.Request.Context().Value(proxy.ReqLogIDKey).(ulid.ULID)
		if !ok {
			return errors.New("reqlog: request is missing ID")
		}

		clone := *res

		// TODO: Use io.LimitReader.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reqlog: could not read response body: %w", err)
		}

		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		clone.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		go func() {
			if err := svc.storeResponse(context.Background(), reqLogID, &clone); err != nil {
				log.Printf("[ERROR] Could not store response log: %v", err)
			}
		}()

		return nil
	}
}
