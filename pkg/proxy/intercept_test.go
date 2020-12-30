package proxy_test

import (
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func getRequests() ([]proxy.RequestEntry, error) {
	return []proxy.RequestEntry{{
		Method: "GET",
		MessageEntry: proxy.MessageEntry{
			UrlEquals: "https://example.com",
			Body:      "<h1>Hello</h1>",
		},
	}}, nil
}

func getResponses() ([]proxy.ResponseEntry, error) {
	return []proxy.ResponseEntry{{
		MessageEntry: proxy.MessageEntry{
			UrlEquals: "https://example.com",
			Headers: map[string]string{
				"Proxy": "Hetty",
			},
		},
	}}, nil
}

// TODO: Come up with more cases

func TestRequestInterceptor(t *testing.T) {
	intercept, _ := proxy.NewIntercept(getRequests, getResponses)
	next := func(req *http.Request) {}

	req, _ := http.NewRequest("GET", "https://example.com", strings.NewReader(""))

	intercept.RequestInterceptor(next)(req)

	body, _ := ioutil.ReadAll(req.Body)

	assert.Equal(t, string(body), "<h1>Hello</h1>")
}

func TestResponseInterceptor(t *testing.T) {
	intercept, _ := proxy.NewIntercept(getRequests, getResponses)
	next := func(res *http.Response) error { return nil }

	req, _ := http.NewRequest("GET", "https://example.com", strings.NewReader(""))

	res := http.Response{
		Request: req,
	}

	intercept.ResponseInterceptor(next)(&res)

	assert.Equal(t, res.Header.Get("Proxy"), "")
}
