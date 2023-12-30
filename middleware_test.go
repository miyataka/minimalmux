package minimalmux

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testMiddleware(s string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(fmt.Sprintf("%s\n", s)))
			h.ServeHTTP(w, r)
		})
	}
}

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("test\n"))
})

func TestAppend(t *testing.T) {
	aFunc := testMiddleware("a")
	bFunc := testMiddleware("b")

	ms := NewMiddlewares()
	ms = ms.Append(aFunc).Append(bFunc)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ms.Handle(testHandler).ServeHTTP(w, r)

	expected := "a\nb\ntest\n"
	if w.Body.String() != expected {
		t.Fatalf("expected %q, got %q", expected, w.Body.String())
	}
}

func TestPrepend(t *testing.T) {
	aFunc := testMiddleware("a")
	bFunc := testMiddleware("b")

	ms := NewMiddlewares()
	ms = ms.Prepend(aFunc).Prepend(bFunc)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	ms.Handle(testHandler).ServeHTTP(w, r)

	expected := "b\na\ntest\n"
	if w.Body.String() != expected {
		t.Fatalf("expected %q, got %q", expected, w.Body.String())
	}
}
