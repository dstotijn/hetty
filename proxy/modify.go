package proxy

import "net/http"

var (
	nopReqModifier = func(req *http.Request) {}
	nopResModifier = func(res *http.Response) error { return nil }
)

type RequestModifyFunc func(req *http.Request)
type ResponseModifyFunc func(res *http.Response) error
