package reqlog_test

//go:generate moq -out proj_mock_test.go -pkg reqlog_test ../proj Service:ProjServiceMock
//go:generate moq -out repo_mock_test.go -pkg reqlog_test . Repository:RepoMock

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dstotijn/hetty/pkg/proj"
	"github.com/dstotijn/hetty/pkg/proxy"
	"github.com/dstotijn/hetty/pkg/reqlog"
)

//nolint:paralleltest
func TestNewService(t *testing.T) {
	projSvcMock := &ProjServiceMock{
		OnProjectOpenFunc:  func(fn proj.OnProjectOpenFn) {},
		OnProjectCloseFunc: func(fn proj.OnProjectCloseFn) {},
	}
	repoMock := &RepoMock{
		FindSettingsByModuleFunc: func(_ context.Context, _ string, _ interface{}) error {
			return nil
		},
	}
	svc := reqlog.NewService(reqlog.Config{
		ProjectService: projSvcMock,
		Repository:     repoMock,
	})

	t.Run("registered handlers for project open and close", func(t *testing.T) {
		got := len(projSvcMock.OnProjectOpenCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.OnProjectOpen` calls (expected: %v, got: %v)", exp, got)
		}

		got = len(projSvcMock.OnProjectCloseCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.OnProjectClose` calls (expected: %v, got: %v)", exp, got)
		}
	})

	t.Run("calls handler when project is opened", func(t *testing.T) {
		// Mock opening a project.
		err := projSvcMock.OnProjectOpenCalls()[0].Fn("foobar")
		if err != nil {
			t.Errorf("unexpected error (expected: nil, got: %v)", err)
		}

		// Assert that settings were fetched from repository, with `svc` as the
		// destination.
		got := len(repoMock.FindSettingsByModuleCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.OnProjectOpen` calls (expected: %v, got: %v)", exp, got)
		}

		findSettingsByModuleCall := repoMock.FindSettingsByModuleCalls()[0]
		expModule := "reqlog"
		expSettings := svc

		if expModule != findSettingsByModuleCall.Module {
			t.Fatalf("incorrect `module` argument for `proj.Service.OnProjectOpen` (expected: %v, got: %v)",
				expModule, findSettingsByModuleCall.Module)
		}

		if expSettings != findSettingsByModuleCall.Settings {
			t.Fatalf("incorrect `settings` argument for `proj.Service.OnProjectOpen` (expected: %v, got: %v)",
				expModule, findSettingsByModuleCall.Settings)
		}
	})

	t.Run("calls handler when project is closed", func(t *testing.T) {
		// Mock updating service settings.
		svc.BypassOutOfScopeRequests = true
		svc.FindReqsFilter = reqlog.FindRequestsFilter{OnlyInScope: true}

		// Mock closing a project.
		err := projSvcMock.OnProjectCloseCalls()[0].Fn("foobar")
		if err != nil {
			t.Errorf("unexpected error (expected: nil, got: %v)", err)
		}

		// Assert that settings were set to defaults on project close.
		expBypassOutOfScopeReqs := false
		expFindReqsFilter := reqlog.FindRequestsFilter{}

		if expBypassOutOfScopeReqs != svc.BypassOutOfScopeRequests {
			t.Fatalf("incorrect `Service.BypassOutOfScopeRequests` value (expected: %v, got: %v)",
				expBypassOutOfScopeReqs, svc.BypassOutOfScopeRequests)
		}

		if expFindReqsFilter != svc.FindReqsFilter {
			t.Fatalf("incorrect `Service.FindReqsFilter` value (expected: %v, got: %v)",
				expFindReqsFilter, svc.FindReqsFilter)
		}
	})
}

//nolint:paralleltest
func TestRequestModifier(t *testing.T) {
	projSvcMock := &ProjServiceMock{
		OnProjectOpenFunc:  func(fn proj.OnProjectOpenFn) {},
		OnProjectCloseFunc: func(fn proj.OnProjectCloseFn) {},
	}
	repoMock := &RepoMock{
		AddRequestLogFunc: func(_ context.Context, _ http.Request, _ []byte, _ time.Time) (*reqlog.Request, error) {
			return &reqlog.Request{}, nil
		},
	}
	svc := reqlog.NewService(reqlog.Config{
		ProjectService: projSvcMock,
		Repository:     repoMock,
	})

	next := func(req *http.Request) {
		req.Body = ioutil.NopCloser(strings.NewReader("modified body"))
	}
	reqModFn := svc.RequestModifier(next)
	req := httptest.NewRequest("GET", "https://example.com/", strings.NewReader("bar"))

	reqModFn(req)

	t.Run("request log was stored in repository", func(t *testing.T) {
		got := len(repoMock.AddRequestLogCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.AddRequestLog` calls (expected: %v, got: %v)", exp, got)
		}
	})

	t.Run("ran next modifier first, before calling repository", func(t *testing.T) {
		got := repoMock.AddRequestLogCalls()[0].Body
		if exp := "modified body"; exp != string(got) {
			t.Fatalf("incorrect `body` argument for `Repository.AddRequestLogCalls` (expected: %v, got: %v)", exp, string(got))
		}
	})
}

//nolint:paralleltest
func TestResponseModifier(t *testing.T) {
	projSvcMock := &ProjServiceMock{
		OnProjectOpenFunc:  func(fn proj.OnProjectOpenFn) {},
		OnProjectCloseFunc: func(fn proj.OnProjectCloseFn) {},
	}
	repoMock := &RepoMock{
		AddResponseLogFunc: func(_ context.Context, _ int64, _ http.Response,
			_ []byte, _ time.Time) (*reqlog.Response, error) {
			return &reqlog.Response{}, nil
		},
	}
	svc := reqlog.NewService(reqlog.Config{
		ProjectService: projSvcMock,
		Repository:     repoMock,
	})

	next := func(res *http.Response) error {
		res.Body = ioutil.NopCloser(strings.NewReader("modified body"))
		return nil
	}
	resModFn := svc.ResponseModifier(next)

	req := httptest.NewRequest("GET", "https://example.com/", strings.NewReader("bar"))
	req = req.WithContext(context.WithValue(req.Context(), proxy.ReqIDKey, int64(42)))

	res := &http.Response{
		Request: req,
		Body:    ioutil.NopCloser(strings.NewReader("bar")),
	}

	if err := resModFn(res); err != nil {
		t.Fatalf("unexpected error (expected: nil, got: %v)", err)
	}

	t.Run("request log was stored in repository", func(t *testing.T) {
		// Dirty (but simple) wait for other goroutine to finish calling repository.
		time.Sleep(10 * time.Millisecond)
		got := len(repoMock.AddResponseLogCalls())
		if exp := 1; exp != got {
			t.Fatalf("incorrect `proj.Service.AddResponseLog` calls (expected: %v, got: %v)", exp, got)
		}
	})

	t.Run("ran next modifier first, before calling repository", func(t *testing.T) {
		got := repoMock.AddResponseLogCalls()[0].Body
		if exp := "modified body"; exp != string(got) {
			t.Fatalf("incorrect `body` argument for `Repository.AddResponseLogCalls` (expected: %v, got: %v)", exp, string(got))
		}
	})
}
