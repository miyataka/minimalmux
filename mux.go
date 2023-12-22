package api

import (
	"net/http"
	"sync"
)

type ServeMux struct {
	tree            *Node
	mu              sync.RWMutex
	notFoundHandler http.HandlerFunc
}

type Route struct {
	Method       string
	Pattern      string
	HandlerFunc  http.HandlerFunc
	PathParamMap map[string]string
}

func (r *Route) IsBlank() bool {
	return r.HandlerFunc == nil
}

func (r *Route) setPathParams(pathParamMap map[string]string) {
	r.PathParamMap = pathParamMap
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		tree:            &Node{},
		notFoundHandler: http.NotFound,
	}
}

const methodAll = "_all"

var methodSlice = []string{
	http.MethodGet,
	http.MethodPost,
	http.MethodPut,
	http.MethodDelete,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodConnect,
	http.MethodTrace,
}

// net/http method wrapper
func (sm *ServeMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sm.handle(methodAll, pattern, handler)
}

func (sm *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	route := sm.tree.search(r.Method, path)
	if route.IsBlank() {
		sm.notFoundHandler(w, r)
		return
	}
	route.HandlerFunc(w, r)
}

func (sm *ServeMux) Handle(pattern string, handler http.Handler) {
	sm.handle(methodAll, pattern, handler.ServeHTTP)
}

func (sm *ServeMux) Handler(r *http.Request) (h http.Handler, pattern string) {
	path := r.URL.Path
	route := sm.tree.search(r.Method, path)
	return route.HandlerFunc, route.Pattern
}

// original method

func (sm *ServeMux) handle(method string, pattern string, handler func(http.ResponseWriter, *http.Request)) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if method == "" {
		panic("http: invalid method")
	}
	if handler == nil {
		panic("http: nil handler")
	}

	// TODO duplicate check?

	if method == methodAll {
		for _, m := range methodSlice {
			sm.tree.insert(m, pattern, Route{
				Method:      m,
				Pattern:     pattern,
				HandlerFunc: handler,
			})
		}
		return
	}

	sm.tree.insert(method, pattern, Route{
		Method:      method,
		Pattern:     pattern,
		HandlerFunc: handler,
	})
}

func (sm *ServeMux) Get(path string, handler http.HandlerFunc) {
	sm.tree.insert(
		http.MethodGet,
		path,
		Route{
			Method:      http.MethodGet,
			Pattern:     path,
			HandlerFunc: handler,
		},
	)
}

func (sm *ServeMux) Post(path string, handler http.HandlerFunc) {
	sm.tree.insert(
		http.MethodPost,
		path,
		Route{
			Method:      http.MethodPost,
			Pattern:     path,
			HandlerFunc: handler,
		},
	)
}

// TODO other http methods
