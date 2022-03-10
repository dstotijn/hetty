package intercept

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"sync"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/log"
	"github.com/dstotijn/hetty/pkg/proxy"
)

var (
	ErrRequestAborted  = errors.New("intercept: request was aborted")
	ErrRequestNotFound = errors.New("intercept: request not found")
	ErrRequestDone     = errors.New("intercept: request is done")
)

// Request represents a server received HTTP request, alongside a channel for sending a modified version of it to the
// routine that's awaiting it. Also contains a channel for receiving a cancellation signal.
type Request struct {
	req  *http.Request
	ch   chan<- *http.Request
	done <-chan struct{}
}

type Service struct {
	mu       *sync.RWMutex
	requests map[ulid.ULID]Request
	logger   log.Logger
}

type Config struct {
	Logger log.Logger
}

// RequestIDs implements sort.Interface.
type RequestIDs []ulid.ULID

func NewService(cfg Config) *Service {
	s := &Service{
		mu:       &sync.RWMutex{},
		requests: make(map[ulid.ULID]Request),
		logger:   cfg.Logger,
	}

	if s.logger == nil {
		s.logger = log.NewNopLogger()
	}

	return s
}

// RequestModifier is a proxy.RequestModifyMiddleware for intercepting HTTP
// requests.
func (svc *Service) RequestModifier(next proxy.RequestModifyFunc) proxy.RequestModifyFunc {
	return func(req *http.Request) {
		// This is a blocking operation, that gets unblocked when either a modified request is returned or an error
		// (typically `context.Canceled`).
		modifiedReq, err := svc.Intercept(req.Context(), req)

		switch {
		case errors.Is(err, ErrRequestAborted):
			svc.logger.Debugw("Stopping intercept, request was aborted.")
			// Prevent further processing by replacing req.Context with a cancelled context value.
			// This will cause the http.Roundtripper in the `proxy` package to
			// handle this request as an error.
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			*req = *req.WithContext(ctx)
		case errors.Is(err, context.Canceled):
			svc.logger.Debugw("Stopping intercept, context was cancelled.")
		case err != nil:
			svc.logger.Errorw("Failed to intercept request.",
				"error", err)
		default:
			*req = *modifiedReq.WithContext(req.Context())
			next(req)
		}
	}
}

// Intercept adds an HTTP request to an array of pending intercepted requests, alongside channels used for sending a
// cancellation signal and receiving a modified request. It's safe for concurrent use.
func (svc *Service) Intercept(ctx context.Context, req *http.Request) (*http.Request, error) {
	reqID, ok := proxy.RequestIDFromContext(ctx)
	if !ok {
		svc.logger.Errorw("Failed to intercept: request doesn't have an ID.")
		return req, nil
	}

	ch := make(chan *http.Request)
	done := make(chan struct{})

	svc.mu.Lock()
	svc.requests[reqID] = Request{
		req:  req,
		ch:   ch,
		done: done,
	}
	svc.mu.Unlock()

	// Whatever happens next (modified request returned, or a context cancelled error), any blocked channel senders
	// should be unblocked, and the request should be removed from the requests queue.
	defer func() {
		close(done)
		svc.mu.Lock()
		defer svc.mu.Unlock()
		delete(svc.requests, reqID)
	}()

	select {
	case modReq := <-ch:
		if modReq == nil {
			return nil, ErrRequestAborted
		}

		return modReq, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ModifyRequest sends a modified HTTP request to the related channel, or returns ErrRequestDone when the request was
// cancelled. It's safe for concurrent use.
func (svc *Service) ModifyRequest(reqID ulid.ULID, modReq *http.Request) error {
	svc.mu.RLock()
	req, ok := svc.requests[reqID]
	svc.mu.RUnlock()

	if !ok {
		return ErrRequestNotFound
	}

	select {
	case <-req.done:
		return ErrRequestDone
	case req.ch <- modReq:
		return nil
	}
}

func (svc *Service) ClearRequests() {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	for _, req := range svc.requests {
		select {
		case <-req.done:
		case req.ch <- nil:
		}
	}
}

// Requests returns a list of pending intercepted requests. It's safe for concurrent use.
func (svc *Service) Requests() []*http.Request {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	ids := make([]ulid.ULID, 0, len(svc.requests))
	for id := range svc.requests {
		ids = append(ids, id)
	}

	sort.Sort(RequestIDs(ids))

	reqs := make([]*http.Request, len(ids))
	for i, id := range ids {
		reqs[i] = svc.requests[id].req
	}

	return reqs
}

// Request returns an intercepted request by ID. It's safe for concurrent use.
func (svc *Service) RequestByID(id ulid.ULID) (*http.Request, error) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	req, ok := svc.requests[id]
	if !ok {
		return nil, ErrRequestNotFound
	}

	return req.req, nil
}

func (ids RequestIDs) Len() int {
	return len(ids)
}

func (ids RequestIDs) Less(i, j int) bool {
	return ids[i].Compare(ids[j]) == -1
}

func (ids RequestIDs) Swap(i, j int) {
	ids[i], ids[j] = ids[j], ids[i]
}
