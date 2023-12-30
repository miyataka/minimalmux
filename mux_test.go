package minimalmux

import (
	"net/http"
	"net/http/httptest"
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
