package reqlog

import (
	"net/http"
	"sync"
)

type Request struct {
	Request http.Request
	Body    []byte
}

type response struct {
	res  http.Response
	body []byte
}

type RequestLogStore struct {
	reqStore []Request
	resStore []response
	reqMu    sync.Mutex
	resMu    sync.Mutex
}

func NewRequestLogStore() RequestLogStore {
	return RequestLogStore{
		reqStore: make([]Request, 0),
		resStore: make([]response, 0),
	}
}

func (store *RequestLogStore) AddRequest(req http.Request, body []byte) {
	store.reqMu.Lock()
	defer store.reqMu.Unlock()

	store.reqStore = append(store.reqStore, Request{req, body})
}

func (store *RequestLogStore) Requests() []Request {
	store.reqMu.Lock()
	defer store.reqMu.Unlock()

	return store.reqStore
}

func (store *RequestLogStore) AddResponse(res http.Response, body []byte) {
	store.resMu.Lock()
	defer store.resMu.Unlock()

	store.resStore = append(store.resStore, response{res, body})
}
