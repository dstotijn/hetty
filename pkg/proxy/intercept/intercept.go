package intercept

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/oklog/ulid"

	"github.com/dstotijn/hetty/pkg/filter"
	"github.com/dstotijn/hetty/pkg/log"
	"github.com/dstotijn/hetty/pkg/proxy"
)

var (
	ErrRequestAborted   = errors.New("intercept: request was aborted")
	ErrRequestNotFound  = errors.New("intercept: request not found")
	ErrRequestDone      = errors.New("intercept: request is done")
	ErrResponseNotFound = errors.New("intercept: response not found")
)

type contextKey int

const interceptResponseKey contextKey = 0

// Request represents a server received HTTP request, alongside a channel for sending a modified version of it to the
// routine that's awaiting it. Also contains a channel for receiving a cancellation signal.
type Request struct {
	req  *http.Request
	ch   chan<- *http.Request
	done <-chan struct{}
}

// Response represents an HTTP response from a proxied request, alongside a channel for sending a modified version of it
// to the routine that's awaiting it. Also contains a channel for receiving a cancellation signal.
type Response struct {
	res  *http.Response
	ch   chan<- *http.Response
	done <-chan struct{}
}

type Item struct {
	Request  *http.Request
	Response *http.Response
}

type Service struct {
	reqMu     *sync.RWMutex
	resMu     *sync.RWMutex
	requests  map[ulid.ULID]Request
	responses map[ulid.ULID]Response
	logger    log.Logger

	requestsEnabled  bool
	responsesEnabled bool
	reqFilter        filter.Expression
	resFilter        filter.Expression
}

type Config struct {
	Logger           log.Logger
	RequestsEnabled  bool
	ResponsesEnabled bool
	RequestFilter    filter.Expression
	ResponseFilter   filter.Expression
}

// RequestIDs implements sort.Interface.
type RequestIDs []ulid.ULID

func NewService(cfg Config) *Service {
	s := &Service{
		reqMu:            &sync.RWMutex{},
		resMu:            &sync.RWMutex{},
		requests:         make(map[ulid.ULID]Request),
		responses:        make(map[ulid.ULID]Response),
		logger:           cfg.Logger,
		requestsEnabled:  cfg.RequestsEnabled,
		responsesEnabled: cfg.ResponsesEnabled,
		reqFilter:        cfg.RequestFilter,
		resFilter:        cfg.ResponseFilter,
	}

	if s.logger == nil {
		s.logger = log.NewNopLogger()
	}

	return s
}

// RequestModifier is a proxy.RequestModifyMiddleware for intercepting HTTP requests.
func (svc *Service) RequestModifier(next proxy.RequestModifyFunc) proxy.RequestModifyFunc {
	return func(req *http.Request) {
		// This is a blocking operation, that gets unblocked when either a modified request is returned or an error
		// (typically `context.Canceled`).
		modifiedReq, err := svc.InterceptRequest(req.Context(), req)

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
			*req = *modifiedReq
			next(req)
		}
	}
}

// InterceptRequest adds an HTTP request to an array of pending intercepted requests, alongside channels used for
// sending a cancellation signal and receiving a modified request. It's safe for concurrent use.
func (svc *Service) InterceptRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	reqID, ok := proxy.RequestIDFromContext(ctx)
	if !ok {
		svc.logger.Errorw("Failed to intercept: context doesn't have an ID.")
		return req, nil
	}

	if !svc.requestsEnabled {
		// If request intercept is disabled, return the incoming request as-is.
		svc.logger.Debugw("Bypassed request interception: feature disabled.")
		return req, nil
	}

	if svc.reqFilter != nil {
		match, err := MatchRequestFilter(req, svc.reqFilter)
		if err != nil {
			return nil, fmt.Errorf("intercept: failed to match request rules for request (id: %v): %w",
				reqID.String(), err,
			)
		}

		if !match {
			svc.logger.Debugw("Bypassed request interception: request rules don't match.")
			return req, nil
		}
	}

	ch := make(chan *http.Request)
	done := make(chan struct{})

	svc.reqMu.Lock()
	svc.requests[reqID] = Request{
		req:  req,
		ch:   ch,
		done: done,
	}
	svc.reqMu.Unlock()

	// Whatever happens next (modified request returned, or a context cancelled error), any blocked channel senders
	// should be unblocked, and the request should be removed from the requests queue.
	defer func() {
		close(done)
		svc.reqMu.Lock()
		defer svc.reqMu.Unlock()
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
func (svc *Service) ModifyRequest(reqID ulid.ULID, modReq *http.Request, modifyResponse *bool) error {
	svc.reqMu.RLock()
	req, ok := svc.requests[reqID]
	svc.reqMu.RUnlock()

	if !ok {
		return ErrRequestNotFound
	}

	*modReq = *modReq.WithContext(req.req.Context())
	if modifyResponse != nil {
		*modReq = *modReq.WithContext(WithInterceptResponse(modReq.Context(), *modifyResponse))
	}

	select {
	case <-req.done:
		return ErrRequestDone
	case req.ch <- modReq:
		return nil
	}
}

// CancelRequest ensures an intercepted request is dropped.
func (svc *Service) CancelRequest(reqID ulid.ULID) error {
	return svc.ModifyRequest(reqID, nil, nil)
}

func (svc *Service) ClearRequests() {
	svc.reqMu.Lock()
	defer svc.reqMu.Unlock()

	for _, req := range svc.requests {
		select {
		case <-req.done:
		case req.ch <- nil:
		}
	}
}

func (svc *Service) ClearResponses() {
	svc.resMu.Lock()
	defer svc.resMu.Unlock()

	for _, res := range svc.responses {
		select {
		case <-res.done:
		case res.ch <- nil:
		}
	}
}

// Items returns a list of pending items (requests and responses). It's safe for concurrent use.
func (svc *Service) Items() []Item {
	svc.reqMu.RLock()
	defer svc.reqMu.RUnlock()

	svc.resMu.RLock()
	defer svc.resMu.RUnlock()

	reqIDs := make([]ulid.ULID, 0, len(svc.requests)+len(svc.responses))

	for id := range svc.requests {
		reqIDs = append(reqIDs, id)
	}

	for id := range svc.responses {
		reqIDs = append(reqIDs, id)
	}

	sort.Sort(RequestIDs(reqIDs))

	items := make([]Item, len(reqIDs))

	for i, id := range reqIDs {
		item := Item{}

		if req, ok := svc.requests[id]; ok {
			item.Request = req.req
		}

		if res, ok := svc.responses[id]; ok {
			item.Response = res.res
		}

		items[i] = item
	}

	return items
}

func (svc *Service) UpdateSettings(settings Settings) {
	// When updating from requests `enabled` -> `disabled`, clear any pending reqs.
	if svc.requestsEnabled && !settings.RequestsEnabled {
		svc.ClearRequests()
	}

	// When updating from responses `enabled` -> `disabled`, clear any pending responses.
	if svc.responsesEnabled && !settings.ResponsesEnabled {
		svc.ClearResponses()
	}

	svc.requestsEnabled = settings.RequestsEnabled
	svc.responsesEnabled = settings.ResponsesEnabled
	svc.reqFilter = settings.RequestFilter
	svc.resFilter = settings.ResponseFilter
}

// ItemByID returns an intercepted item (request and possible response) by ID. It's safe for concurrent use.
func (svc *Service) ItemByID(id ulid.ULID) (Item, error) {
	svc.reqMu.RLock()
	defer svc.reqMu.RUnlock()

	svc.resMu.RLock()
	defer svc.resMu.RUnlock()

	item := Item{}
	found := false

	if req, ok := svc.requests[id]; ok {
		item.Request = req.req
		found = true
	}

	if res, ok := svc.responses[id]; ok {
		item.Response = res.res
		found = true
	}

	if !found {
		return Item{}, ErrRequestNotFound
	}

	return item, nil
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

func WithInterceptResponse(ctx context.Context, value bool) context.Context {
	return context.WithValue(ctx, interceptResponseKey, value)
}

func ShouldInterceptResponseFromContext(ctx context.Context) (bool, bool) {
	shouldIntercept, ok := ctx.Value(interceptResponseKey).(bool)
	return shouldIntercept, ok
}

// ResponseModifier is a proxy.ResponseModifyMiddleware for intercepting HTTP responses.
func (svc *Service) ResponseModifier(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
	return func(res *http.Response) error {
		// This is a blocking operation, that gets unblocked when either a modified response is returned or an error.
		//nolint:bodyclose
		modifiedRes, err := svc.InterceptResponse(res.Request.Context(), res)
		if err != nil {
			return fmt.Errorf("failed to intercept response: %w", err)
		}

		*res = *modifiedRes

		return next(res)
	}
}

// InterceptResponse adds an HTTP response to an array of pending intercepted responses, alongside channels used for
// sending a cancellation signal and receiving a modified response. It's safe for concurrent use.
func (svc *Service) InterceptResponse(ctx context.Context, res *http.Response) (*http.Response, error) {
	reqID, ok := proxy.RequestIDFromContext(ctx)
	if !ok {
		svc.logger.Errorw("Failed to intercept: context doesn't have an ID.")
		return res, nil
	}

	shouldIntercept, ok := ShouldInterceptResponseFromContext(ctx)
	if ok && !shouldIntercept {
		// If the related request explicitly disabled response intercept, return the response as-is.
		svc.logger.Debugw("Bypassed response interception: related request explicitly disabled response intercept.")
		return res, nil
	}

	// If global response intercept is disabled and interception is *not* explicitly enabled for this response: bypass.
	if !svc.responsesEnabled && !(ok && shouldIntercept) {
		svc.logger.Debugw("Bypassed response interception: feature disabled.")
		return res, nil
	}

	if svc.resFilter != nil {
		match, err := MatchResponseFilter(res, svc.resFilter)
		if err != nil {
			return nil, fmt.Errorf("intercept: failed to match response rules for response (id: %v): %w",
				reqID.String(), err,
			)
		}

		if !match {
			svc.logger.Debugw("Bypassed response interception: response rules don't match.")
			return res, nil
		}
	}

	ch := make(chan *http.Response)
	done := make(chan struct{})

	svc.resMu.Lock()
	svc.responses[reqID] = Response{
		res:  res,
		ch:   ch,
		done: done,
	}
	svc.resMu.Unlock()

	// Whatever happens next (modified response returned, or a context cancelled error), any blocked channel senders
	// should be unblocked, and the response should be removed from the responses queue.
	defer func() {
		close(done)
		svc.resMu.Lock()
		defer svc.resMu.Unlock()
		delete(svc.responses, reqID)
	}()

	select {
	case modRes := <-ch:
		if modRes == nil {
			return nil, ErrRequestAborted
		}

		return modRes, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// ModifyResponse sends a modified HTTP response to the related channel, or returns ErrRequestDone when the related
// request was cancelled. It's safe for concurrent use.
func (svc *Service) ModifyResponse(reqID ulid.ULID, modRes *http.Response) error {
	svc.resMu.RLock()
	res, ok := svc.responses[reqID]
	svc.resMu.RUnlock()

	if !ok {
		return ErrRequestNotFound
	}

	if modRes != nil {
		modRes.Request = res.res.Request
	}

	select {
	case <-res.done:
		return ErrRequestDone
	case res.ch <- modRes:
		return nil
	}
}

// CancelResponse ensures an intercepted response is dropped.
func (svc *Service) CancelResponse(reqID ulid.ULID) error {
	return svc.ModifyResponse(reqID, nil)
}
