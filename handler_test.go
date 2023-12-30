package minimalmux

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthCheckHandler(t *testing.T) {
	mux := NewServeMux()
	mux.HandleFunc("/health", HealthcheckHandler)

	ts := httptest.NewServer(mux)

	res, err := http.DefaultClient.Get(ts.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}
}
