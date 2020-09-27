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
	RequestID uuid.UUID
	Response  http.Response
	Body      []byte
	Timestamp time.Time
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo}
}

func (svc *Service) FindAllRequests(ctx context.Context) ([]Request, error) {
	return svc.repo.FindAllRequestLogs(ctx)
}

func (svc *Service) FindRequestLogByID(ctx context.Context, id uuid.UUID) (Request, error) {
	return svc.repo.FindRequestLogByID(ctx, id)
}

func (svc *Service) addRequest(ctx context.Context, reqID uuid.UUID, req http.Request, body []byte) error {
	reqLog := Request{
		ID:        reqID,
		Request:   req,
		Body:      body,
		Timestamp: time.Now(),
	}

	return svc.repo.AddRequestLog(ctx, reqLog)
}

func (svc *Service) addResponse(ctx context.Context, reqID uuid.UUID, res http.Response, body []byte) error {
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

	resLog := Response{
		RequestID: reqID,
		Response:  res,
		Body:      body,
		Timestamp: time.Now(),
	}

	return svc.repo.AddResponseLog(ctx, resLog)
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

		go func() {
			if err := svc.addRequest(context.Background(), reqID, *clone, body); err != nil {
				log.Printf("[ERROR] Could not store request log: %v", err)
			}
		}()
	}
}

func (svc *Service) ResponseModifier(next proxy.ResponseModifyFunc) proxy.ResponseModifyFunc {
	return func(res *http.Response) error {
		if err := next(res); err != nil {
			return err
		}

		reqID, _ := res.Request.Context().Value(proxy.ReqIDKey).(uuid.UUID)
		if reqID == uuid.Nil {
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
			if err := svc.addResponse(res.Request.Context(), reqID, clone, body); err != nil {
				log.Printf("[ERROR] Could not store response log: %v", err)
			}
		}()

		return nil
	}
}
