package api

import "net/http"

type ServeMux struct {
	mux    *http.ServeMux
	routes []Route
}

type Route struct {
	Method  string
	Path    string
	Handler http.Handler
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		mux:    http.NewServeMux(), // TODO
		routes: []Route{},
	}
}

// net/http method wrapper
func (sm *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sm.mux.HandleFunc(pattern, handler)
}

func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
	sm.mux.ServeHTTP(w, r)
}

func (sm *ServeMux) Handle(pattern string, handler http.Handler) {
	sm.mux.Handle(pattern, handler)
}

func (sm *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return sm.mux.Handler(r)
}

// original method
func (sm *ServeMux) Get(path string, handler http.Handler) {
	sm.routes = append(sm.routes, Route{
		Method:  http.MethodGet,
		Path:    path,
		Handler: handler,
	})
}

func (sm *ServeMux) Post(path string, handler http.Handler) {
	sm.routes = append(sm.routes, Route{
		Method:  http.MethodPost,
		Path:    path,
		Handler: handler,
	})
}

// TODO other http methods
