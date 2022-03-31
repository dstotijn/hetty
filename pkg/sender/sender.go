package sender

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/filter"
	"github.com/dstotijn/hetty/pkg/reqlog"
	"github.com/dstotijn/hetty/pkg/scope"
)

//nolint:gosec
var ulidEntropy = rand.New(rand.NewSource(time.Now().UnixNano()))

var defaultHTTPClient = &http.Client{
	Transport: &HTTPTransport{},
	Timeout:   30 * time.Second,
}

var (
	ErrProjectIDMustBeSet = errors.New("sender: project ID must be set")
	ErrRequestNotFound    = errors.New("sender: request not found")
)

type Service interface {
	FindRequestByID(ctx context.Context, id ulid.ULID) (Request, error)
	FindRequests(ctx context.Context) ([]Request, error)
	CreateOrUpdateRequest(ctx context.Context, req Request) (Request, error)
	CloneFromRequestLog(ctx context.Context, reqLogID ulid.ULID) (Request, error)
	DeleteRequests(ctx context.Context, projectID ulid.ULID) error
	SendRequest(ctx context.Context, id ulid.ULID) (Request, error)
	SetActiveProjectID(ulid.ULID)
	SetFindReqsFilter(filter FindRequestsFilter)
	FindReqsFilter() FindRequestsFilter
}

type service struct {
	activeProjectID ulid.ULID
	findReqsFilter  FindRequestsFilter
	scope           *scope.Scope
	repo            Repository
	reqLogSvc       reqlog.Service
	httpClient      *http.Client
}

type FindRequestsFilter struct {
	ProjectID   ulid.ULID
	OnlyInScope bool
	SearchExpr  filter.Expression
}

type Config struct {
	Scope         *scope.Scope
	Repository    Repository
	ReqLogService reqlog.Service
	HTTPClient    *http.Client
}

type SendError struct {
	err error
}

func NewService(cfg Config) Service {
	svc := &service{
		repo:       cfg.Repository,
		reqLogSvc:  cfg.ReqLogService,
		httpClient: defaultHTTPClient,
		scope:      cfg.Scope,
	}

	if cfg.HTTPClient != nil {
		svc.httpClient = cfg.HTTPClient
	}

	return svc
}

type Request struct {
	ID                 ulid.ULID
	ProjectID          ulid.ULID
	SourceRequestLogID ulid.ULID

	URL    *url.URL
	Method string
	Proto  string
	Header http.Header
	Body   []byte

	Response *reqlog.ResponseLog
}

func (svc *service) FindRequestByID(ctx context.Context, id ulid.ULID) (Request, error) {
	req, err := svc.repo.FindSenderRequestByID(ctx, id)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to find request: %w", err)
	}

	return req, nil
}

func (svc *service) FindRequests(ctx context.Context) ([]Request, error) {
	return svc.repo.FindSenderRequests(ctx, svc.findReqsFilter, svc.scope)
}

func (svc *service) CreateOrUpdateRequest(ctx context.Context, req Request) (Request, error) {
	if svc.activeProjectID.Compare(ulid.ULID{}) == 0 {
		return Request{}, ErrProjectIDMustBeSet
	}

	if req.ID.Compare(ulid.ULID{}) == 0 {
		req.ID = ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy)
	}

	req.ProjectID = svc.activeProjectID

	if req.Method == "" {
		req.Method = http.MethodGet
	}

	if req.Proto == "" {
		req.Proto = HTTPProto20
	}

	if !isValidProto(req.Proto) {
		return Request{}, fmt.Errorf("sender: unsupported HTTP protocol: %v", req.Proto)
	}

	err := svc.repo.StoreSenderRequest(ctx, req)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to store request: %w", err)
	}

	return req, nil
}

func (svc *service) CloneFromRequestLog(ctx context.Context, reqLogID ulid.ULID) (Request, error) {
	if svc.activeProjectID.Compare(ulid.ULID{}) == 0 {
		return Request{}, ErrProjectIDMustBeSet
	}

	reqLog, err := svc.reqLogSvc.FindRequestLogByID(ctx, reqLogID)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to find request log: %w", err)
	}

	req := Request{
		ID:                 ulid.MustNew(ulid.Timestamp(time.Now()), ulidEntropy),
		ProjectID:          svc.activeProjectID,
		SourceRequestLogID: reqLogID,
		Method:             reqLog.Method,
		URL:                reqLog.URL,
		Proto:              HTTPProto20, // Attempt HTTP/2.
		Header:             reqLog.Header,
		Body:               reqLog.Body,
	}

	err = svc.repo.StoreSenderRequest(ctx, req)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to store request: %w", err)
	}

	return req, nil
}

func (svc *service) SetFindReqsFilter(filter FindRequestsFilter) {
	svc.findReqsFilter = filter
}

func (svc *service) FindReqsFilter() FindRequestsFilter {
	return svc.findReqsFilter
}

func (svc *service) SendRequest(ctx context.Context, id ulid.ULID) (Request, error) {
	req, err := svc.repo.FindSenderRequestByID(ctx, id)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to find request: %w", err)
	}

	httpReq, err := parseHTTPRequest(ctx, req)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to parse HTTP request: %w", err)
	}

	resLog, err := svc.sendHTTPRequest(httpReq)
	if err != nil {
		return Request{}, fmt.Errorf("sender: could not send HTTP request: %w", err)
	}

	err = svc.repo.StoreResponseLog(ctx, id, resLog)
	if err != nil {
		return Request{}, fmt.Errorf("sender: failed to store sender response log: %w", err)
	}

	req.Response = &resLog

	return req, nil
}

func parseHTTPRequest(ctx context.Context, req Request) (*http.Request, error) {
	ctx = context.WithValue(ctx, protoCtxKey{}, req.Proto)

	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL.String(), bytes.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to construct HTTP request: %w", err)
	}

	if req.Header != nil {
		httpReq.Header = req.Header
	}

	return httpReq, nil
}

func (svc *service) sendHTTPRequest(httpReq *http.Request) (reqlog.ResponseLog, error) {
	res, err := svc.httpClient.Do(httpReq)
	if err != nil {
		return reqlog.ResponseLog{}, &SendError{err}
	}
	defer res.Body.Close()

	resLog, err := reqlog.ParseHTTPResponse(res)
	if err != nil {
		return reqlog.ResponseLog{}, fmt.Errorf("failed to parse http response: %w", err)
	}

	return resLog, err
}

func (svc *service) SetActiveProjectID(id ulid.ULID) {
	svc.activeProjectID = id
}

func (svc *service) DeleteRequests(ctx context.Context, projectID ulid.ULID) error {
	return svc.repo.DeleteSenderRequests(ctx, projectID)
}

func (e SendError) Error() string {
	return fmt.Sprintf("failed to send HTTP request: %v", e.err)
}

func (e SendError) Unwrap() error {
	return e.err
}
