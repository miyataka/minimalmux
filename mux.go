package minimalmux

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
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

	if method == methodAll {
		for _, m := range methodSlice {
			// duplicate check
			r := sm.tree.search(method, pattern)
			if r.Method == method && r.Pattern == pattern {
				panic("http: multiple registrations for " + method + " " + pattern)
			}

			// insert
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
	sm.method(http.MethodGet, path, handler)
}

func (sm *ServeMux) Post(path string, handler http.HandlerFunc) {
	sm.method(http.MethodPost, path, handler)
}

func (sm *ServeMux) Put(path string, handler http.HandlerFunc) {
	sm.method(http.MethodPut, path, handler)
}

func (sm *ServeMux) Delete(path string, handler http.HandlerFunc) {
	sm.method(http.MethodDelete, path, handler)
}

func (sm *ServeMux) Head(path string, handler http.HandlerFunc) {
	sm.method(http.MethodHead, path, handler)
}

func (sm *ServeMux) Options(path string, handler http.HandlerFunc) {
	sm.method(http.MethodOptions, path, handler)
}

func (sm *ServeMux) Patch(path string, handler http.HandlerFunc) {
	sm.method(http.MethodPatch, path, handler)
}

func (sm *ServeMux) method(method string, path string, handler http.HandlerFunc) {
	sm.tree.insert(
		method,
		path,
		Route{
			Method:      method,
			Pattern:     path,
			HandlerFunc: handler,
		},
	)
}

// graceful shutdown
// TODO use custom type instead of  time.Duration
func ListenAndServeWithGracefulShutdown(srv *http.Server, d time.Duration) error {
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGTERM, os.Interrupt, os.Kill,
	)
	defer stop()

	ctx, cancelCauseFunc := context.WithCancelCause(ctx)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			cancelCauseFunc(err)
		}
	}()

	<-ctx.Done() // wait signal

	// check error caused by cancel or not
	if err := context.Cause(ctx); err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}

	// shutdown

	ctx, cancelFunc := context.WithTimeout(context.Background(), d)
	defer cancelFunc()

	// shutdown server with timeout
	var shutdownErr error = nil
	if err := srv.Shutdown(ctx); err != nil {
		shutdownErr = err
	}

	// check timeout occurred or not
	if err := context.Cause(ctx); err != nil || shutdownErr != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errors.Join(err, shutdownErr)
		}
		return shutdownErr
	}

	return nil
}
