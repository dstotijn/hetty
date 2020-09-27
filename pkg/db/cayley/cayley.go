package cayley

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/cayleygraph/cayley"
	"github.com/cayleygraph/cayley/graph"
	"github.com/cayleygraph/cayley/graph/kv"
	"github.com/cayleygraph/cayley/schema"
	"github.com/cayleygraph/quad"
	"github.com/cayleygraph/quad/voc"
	"github.com/cayleygraph/quad/voc/rdf"
	"github.com/google/uuid"

	"github.com/dstotijn/hetty/pkg/reqlog"
)

type HTTPRequest struct {
	rdfType   struct{}      `quad:"@type > hy:HTTPRequest"`
	ID        quad.IRI      `quad:"@id"`
	Proto     string        `quad:"hy:proto"`
	URL       string        `quad:"hy:url"`
	Method    string        `quad:"hy:method"`
	Body      string        `quad:"hy:body,optional"`
	Headers   []HTTPHeader  `quad:"hy:header"`
	Timestamp time.Time     `quad:"hy:timestamp"`
	Response  *HTTPResponse `quad:"hy:request < *,optional"`
}

type HTTPResponse struct {
	rdfType    struct{}     `quad:"@type > hy:HTTPResponse"`
	RequestID  quad.IRI     `quad:"hy:request"`
	Proto      string       `quad:"hy:proto"`
	Status     string       `quad:"hy:status"`
	StatusCode int          `quad:"hy:status_code"`
	Headers    []HTTPHeader `quad:"hy:header"`
	Body       string       `quad:"hy:body,optional"`
	Timestamp  time.Time    `quad:"hy:timestamp"`
}

type HTTPHeader struct {
	rdfType struct{} `quad:"@type > hy:HTTPHeader"`
	Key     string   `quad:"hy:key"`
	Value   string   `quad:"hy:value,optional"`
}

type Database struct {
	store  *cayley.Handle
	schema *schema.Config
	mu     sync.Mutex
}

func init() {
	voc.RegisterPrefix("hy:", "https://hetty.xyz/")
	schema.RegisterType(quad.IRI("hy:HTTPRequest"), HTTPRequest{})
	schema.RegisterType(quad.IRI("hy:HTTPResponse"), HTTPResponse{})
	schema.RegisterType(quad.IRI("hy:HTTPHeader"), HTTPHeader{})

	kv.Register(Type, kv.Registration{
		NewFunc:      boltOpen,
		InitFunc:     boltCreate,
		IsPersistent: true,
	})
}

func NewDatabase(filename string) (*Database, error) {
	dir, file := path.Split(filename)
	if dir == "" {
		dir = "."
	}
	opts := graph.Options{
		"filename": file,
	}

	schemaCfg := schema.NewConfig()
	schemaCfg.GenerateID = func(_ interface{}) quad.Value {
		return quad.BNode(uuid.New().String())
	}

	// Initialize the database.
	err := graph.InitQuadStore("bolt", dir, opts)
	if err != nil && err != graph.ErrDatabaseExists {
		return nil, fmt.Errorf("cayley: could not initialize database: %v", err)
	}

	// Open the database.
	store, err := cayley.NewGraph("bolt", dir, opts)
	if err != nil {
		return nil, fmt.Errorf("cayley: could not open database: %v", err)
	}

	return &Database{
		store:  store,
		schema: schemaCfg,
	}, nil
}

func (db *Database) Close() error {
	return db.store.Close()
}

func (db *Database) FindAllRequestLogs(ctx context.Context) ([]reqlog.Request, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var reqLogs []reqlog.Request
	var reqs []HTTPRequest

	path := cayley.StartPath(db.store, quad.IRI("hy:HTTPRequest")).In(quad.IRI(rdf.Type))
	err := path.Iterate(ctx).EachValue(db.store, func(v quad.Value) {
		var req HTTPRequest
		if err := db.schema.LoadToDepth(ctx, db.store, &req, -1, v); err != nil {
			log.Printf("[ERROR] Could not load sub-graph for http requests: %v", err)
			return
		}
		reqs = append(reqs, req)
	})
	if err != nil {
		return nil, fmt.Errorf("cayley: could not iterate over http requests: %v", err)
	}

	for _, req := range reqs {
		reqLog, err := parseRequestQuads(req, nil)
		if err != nil {
			return nil, fmt.Errorf("cayley: could not parse request quads (id: %v): %v", req.ID, err)
		}
		reqLogs = append(reqLogs, reqLog)
	}

	// By default, all retrieved requests are ordered chronologically, oldest first.
	// Reverse the order, so newest logs are first.
	for i := len(reqLogs)/2 - 1; i >= 0; i-- {
		opp := len(reqLogs) - 1 - i
		reqLogs[i], reqLogs[opp] = reqLogs[opp], reqLogs[i]
	}

	return reqLogs, nil
}

func (db *Database) FindRequestLogByID(ctx context.Context, id uuid.UUID) (reqlog.Request, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	var req HTTPRequest
	err := db.schema.LoadTo(ctx, db.store, &req, iriFromUUID(id))
	if schema.IsNotFound(err) {
		return reqlog.Request{}, reqlog.ErrRequestNotFound
	}
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("cayley: could not load value: %v", err)
	}

	reqLog, err := parseRequestQuads(req, nil)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("cayley: could not parse request log (id: %v): %v", req.ID, err)
	}

	return reqLog, nil
}

func (db *Database) AddRequestLog(ctx context.Context, reqLog reqlog.Request) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	httpReq := HTTPRequest{
		ID:        iriFromUUID(reqLog.ID),
		Proto:     reqLog.Request.Proto,
		Method:    reqLog.Request.Method,
		URL:       reqLog.Request.URL.String(),
		Headers:   httpHeadersSliceFromMap(reqLog.Request.Header),
		Body:      string(reqLog.Body),
		Timestamp: reqLog.Timestamp,
	}

	tx := cayley.NewTransaction()
	qw := graph.NewTxWriter(tx, graph.Add)

	_, err := db.schema.WriteAsQuads(qw, httpReq)
	if err != nil {
		return fmt.Errorf("cayley: could not write quads: %v", err)
	}

	if err := db.store.ApplyTransaction(tx); err != nil {
		return fmt.Errorf("cayley: could not apply transaction: %v", err)
	}

	return nil
}

func (db *Database) AddResponseLog(ctx context.Context, resLog reqlog.Response) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	httpRes := HTTPResponse{
		RequestID:  iriFromUUID(resLog.RequestID),
		Proto:      resLog.Response.Proto,
		Status:     resLog.Response.Status,
		StatusCode: resLog.Response.StatusCode,
		Headers:    httpHeadersSliceFromMap(resLog.Response.Header),
		Body:       string(resLog.Body),
		Timestamp:  resLog.Timestamp,
	}

	tx := cayley.NewTransaction()
	qw := graph.NewTxWriter(tx, graph.Add)

	_, err := db.schema.WriteAsQuads(qw, httpRes)
	if err != nil {
		return fmt.Errorf("cayley: could not write response quads: %v", err)
	}

	if err := db.store.ApplyTransaction(tx); err != nil {
		return fmt.Errorf("cayley: could not apply transaction: %v", err)
	}

	return nil
}

func iriFromUUID(id uuid.UUID) quad.IRI {
	return quad.IRI("hy:" + id.String()).Full().Short()
}

func uuidFromIRI(iri quad.IRI) (uuid.UUID, error) {
	iriString := iri.Short().String()
	stripped := strings.TrimRight(strings.TrimLeft(iriString, "<hy:"), ">")
	id, err := uuid.Parse(stripped)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func httpHeadersSliceFromMap(hm http.Header) []HTTPHeader {
	if hm == nil {
		return nil
	}
	var hs []HTTPHeader
	for key, values := range hm {
		for _, value := range values {
			hs = append(hs, HTTPHeader{Key: key, Value: value})
		}
	}
	return hs
}

func httpHeadersMapFromSlice(hs []HTTPHeader) http.Header {
	if hs == nil {
		return nil
	}
	hm := make(http.Header)
	for _, header := range hs {
		hm.Add(header.Key, header.Value)
	}
	return hm
}

func parseRequestQuads(req HTTPRequest, _ *HTTPResponse) (reqlog.Request, error) {
	reqID, err := uuidFromIRI(req.ID)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("cannot parse request id: %v", err)
	}

	u, err := url.Parse(req.URL)
	if err != nil {
		return reqlog.Request{}, fmt.Errorf("cannot parse request url: %v", err)
	}

	reqLog := reqlog.Request{
		ID: reqID,
		Request: http.Request{
			Method: req.Method,
			URL:    u,
			Proto:  req.Proto,
			Header: httpHeadersMapFromSlice(req.Headers),
		},
		Timestamp: req.Timestamp,
	}
	if req.Body != "" {
		reqLog.Body = []byte(reqLog.Body)
	}

	if req.Response != nil {
		reqLog.Response = &reqlog.Response{
			RequestID: reqID,
			Response: http.Response{
				Proto:      req.Response.Proto,
				Status:     req.Response.Status,
				StatusCode: req.Response.StatusCode,
				Header:     httpHeadersMapFromSlice(req.Response.Headers),
			},
		}
		if req.Response.Body != "" {
			reqLog.Response.Body = []byte(req.Response.Body)
		}
	}

	return reqLog, nil
}
