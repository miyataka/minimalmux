package minimalmux

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestServeMuxRouting(t *testing.T) {
	tcs := []struct {
		method       string
		pattern      string
		expectStatus int
	}{
		{method: http.MethodGet, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodPost, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodPut, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodDelete, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodHead, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodOptions, pattern: "/test", expectStatus: http.StatusOK},
		{method: http.MethodPatch, pattern: "/test", expectStatus: http.StatusOK},
	}

	for _, tc := range tcs {
		mux := NewServeMux()
		mux.handle(tc.method, tc.pattern, testHandler)

		testServer := httptest.NewServer(mux)
		res := testHttpRequest(t, testServer, tc.method, tc.pattern)

		if res.StatusCode != tc.expectStatus {
			t.Errorf("Status code not equal. got: %d, want: %d", res.StatusCode, tc.expectStatus)
		}
	}
}

func testHttpRequest(t *testing.T, ts *httptest.Server, method, path string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

func TestUnregisterRoute(t *testing.T) {
	registeredRoutes := []struct {
		method       string
		pattern      string
		expectStatus int
	}{
		{
			method:       http.MethodGet,
			pattern:      "/test",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodPost,
			pattern:      "/test/foo",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodPut,
			pattern:      "/test/foo/:id",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodDelete,
			pattern:      "/test/foo/:id",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodDelete,
			pattern:      "/test/foo/:id",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodPost,
			pattern:      "/test/foo/:id/bar",
			expectStatus: http.StatusOK,
		},
		{
			method:       http.MethodPut,
			pattern:      "/test/foo/:id/bar/:barID",
			expectStatus: http.StatusOK,
		},
	}

	mux := NewServeMux()
	for _, r := range registeredRoutes {
		mux.handle(r.method, r.pattern, testHandler)
	}

	ts := httptest.NewServer(mux)

	tcs := []struct {
		req          http.Request
		expectStatus int
	}{
		{
			req:          http.Request{Method: http.MethodGet, URL: &url.URL{Path: "/test"}},
			expectStatus: http.StatusOK,
		},
		{ // FIXME: should return `method not allowed`?
			req:          http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/test"}},
			expectStatus: http.StatusNotFound,
		},
		{
			req:          http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/test/foo"}},
			expectStatus: http.StatusOK,
		},
		{
			req:          http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/test/foo/1"}},
			expectStatus: http.StatusNotFound,
		},
		{
			req:          http.Request{Method: http.MethodPut, URL: &url.URL{Path: "/test/foo/this-is-id"}},
			expectStatus: http.StatusOK,
		},
		{
			req:          http.Request{Method: http.MethodDelete, URL: &url.URL{Path: "/test/foo/this-is-id"}},
			expectStatus: http.StatusOK,
		},
		{
			req:          http.Request{Method: http.MethodPost, URL: &url.URL{Path: "/test/foo/123/bar"}},
			expectStatus: http.StatusOK,
		},
		{
			req:          http.Request{Method: http.MethodPut, URL: &url.URL{Path: "/test/foo/123/bar/456"}},
			expectStatus: http.StatusOK,
		},
		{ // FIXME: tailing slash should be ignored??
			req:          http.Request{Method: http.MethodPut, URL: &url.URL{Path: "/test/foo/123/bar/"}},
			expectStatus: http.StatusOK,
		},
	}

	for _, tc := range tcs {
		res := testHttpRequest(t, ts, tc.req.Method, tc.req.URL.Path)
		if res.StatusCode != tc.expectStatus {
			t.Errorf("Status code not equal. got: %d, want: %d", res.StatusCode, tc.expectStatus)
		}
	}
}
