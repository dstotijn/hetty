package proxy

import "net/http"

var (
	nopReqModifier = func(req *http.Request) {}
	nopResModifier = func(res *http.Response) error { return nil }
)

// RequestModifyFunc defines a type for a function that can modify a HTTP
// request before it's proxied.
type RequestModifyFunc func(req *http.Request)

// RequestModifyMiddleware defines a type for chaining request modifier
// middleware.
type RequestModifyMiddleware func(next RequestModifyFunc) RequestModifyFunc

// ResponseModifyFunc defines a type for a function that can modify a HTTP
// response before it's written back to the client.
type ResponseModifyFunc func(res *http.Response) error

// ResponseModifyMiddleware defines a type for chaining response modifier
// middleware.
type ResponseModifyMiddleware func(ResponseModifyFunc) ResponseModifyFunc
