package reqlog

import (
	"net/http"
	"sync"
)

type request struct {
	req  http.Request
	body []byte
}

type response struct {
	res  http.Response
	body []byte
}

type RequestLog struct {
	reqStore []request
	resStore []response
	reqMu    sync.Mutex
	resMu    sync.Mutex
}

func NewRequestLog() RequestLog {
	return RequestLog{
		reqStore: make([]request, 0),
		resStore: make([]response, 0),
	}
}

func (rl *RequestLog) AddRequest(req http.Request, body []byte) {
	rl.reqMu.Lock()
	defer rl.reqMu.Unlock()

	rl.reqStore = append(rl.reqStore, request{req, body})
}

func (rl *RequestLog) AddResponse(res http.Response, body []byte) {
	rl.resMu.Lock()
	defer rl.resMu.Unlock()

	rl.resStore = append(rl.resStore, response{res, body})
}
