package api

import "net/http"

type ServeMux struct {
	mux    *http.ServeMux
	routes []Route
	tree   *Node
}

type Route struct {
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		mux:    http.NewServeMux(), // TODO
		routes: []Route{},
		tree:   &Node{},
	}
}

// net/http method wrapper
func (sm *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sm.mux.HandleFunc(pattern, handler)
}

func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	route := sm.tree.search(path)
	if route.HandlerFunc != nil {
		route.HandlerFunc(w, r)
	} else {
		sm.mux.ServeHTTP(w, r)
	}
}

func (sm *ServeMux) Handle(pattern string, handler http.Handler) {
	sm.mux.Handle(pattern, handler)
}

func (sm *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	return sm.mux.Handler(r)
}

// original method
func (sm *ServeMux) setupRouteTree() {
	// TODO method routing
	// TODO host/domain routing
	for _, r := range sm.routes {
		sm.tree.insert(r.Path, r)
	}
}

func (sm *ServeMux) Get(path string, handler http.Handler) {
	sm.routes = append(sm.routes, Route{
		Method:      http.MethodGet,
		Path:        path,
		HandlerFunc: handler.ServeHTTP,
	})
}

func (sm *ServeMux) Post(path string, handler http.Handler) {
	sm.routes = append(sm.routes, Route{
		Method:      http.MethodPost,
		Path:        path,
		HandlerFunc: handler.ServeHTTP,
	})
}

// TODO other http methods
