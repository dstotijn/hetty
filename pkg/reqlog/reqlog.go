package reqlog

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/google/uuid"
)

var ErrRequestNotFound = errors.New("reqlog: request not found")

type Request struct {
	ID        uuid.UUID
	Request   http.Request
	Body      []byte
	Timestamp time.Time
	Response  *Response
}

type Response struct {
	Response  http.Response
	Body      []byte
	Timestamp time.Time
}

type Service struct {
	store []Request
	mu    sync.Mutex
}

func NewService() Service {
	return Service{
		store: make([]Request, 0),
	}
}

func (svc *Service) FindAllRequests() []Request {
	return svc.store
}

func (svc *Service) FindRequestLogByID(id string) (Request, error) {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	for _, req := range svc.store {
		if req.ID.String() == id {
			return req, nil
		}
	}

	return Request{}, ErrRequestNotFound
}

func (svc *Service) addRequest(reqID uuid.UUID, req http.Request, body []byte) Request {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	reqLog := Request{
		ID:        reqID,
		Request:   req,
		Body:      body,
		Timestamp: time.Now(),
	}

	svc.store = append(svc.store, reqLog)

	return reqLog
}

func (svc *Service) addResponse(reqID uuid.UUID, res http.Response, body []byte) error {
	svc.mu.Lock()
	defer svc.mu.Unlock()

	for i := range svc.store {
		if svc.store[i].ID == reqID {
			svc.store[i].Response = &Response{
				Response:  res,
				Body:      body,
				Timestamp: time.Now(),
			}
			return nil
		}
	}

	return fmt.Errorf("no request found with ID: %s", reqID)
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
		}

		reqID, _ := req.Context().Value(proxy.ReqIDKey).(uuid.UUID)
		if reqID == uuid.Nil {
			log.Println("[ERROR] Request is missing a related request ID")
			return
		}

		_ = svc.addRequest(reqID, *clone, body)
	}
}

func (svc *Service) ResponseModifier(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		clone := *res

		// TODO: Use io.LimitReader.
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("reqlog: could not read response body: %v", err)
		}
		res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

		if res.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, err := gzip.NewReader(bytes.NewBuffer(body))
			if err != nil {
				return fmt.Errorf("reqlog: could not create gzip reader: %v", err)
			}
			defer gzipReader.Close()
			body, err = ioutil.ReadAll(gzipReader)
			if err != nil {
				return fmt.Errorf("reqlog: could not read gzipped response body: %v", err)
			}
		}

		reqID, _ := res.Request.Context().Value(proxy.ReqIDKey).(uuid.UUID)
		if reqID == uuid.Nil {
			return errors.New("reqlog: request is missing ID")
		}

		if err := svc.addResponse(reqID, clone, body); err != nil {
			return fmt.Errorf("reqlog: could not add response: %v", err)
		}

		return nil
	}
}
