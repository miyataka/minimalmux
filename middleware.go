package minimalmux

import "net/http"

type Middlewares []func(http.Handler) http.Handler

func NewMiddlewares() *Middlewares {
	return &Middlewares{}
}

func (m *Middlewares) Append(middleware func(http.Handler) http.Handler) *Middlewares {
	*m = append(*m, middleware)
	return m
}

func (m *Middlewares) Prepend(middleware func(http.Handler) http.Handler) *Middlewares {
	*m = append(Middlewares{middleware}, *m...)
	return m
}

func (m *Middlewares) Handle(h http.Handler) http.Handler {
	if len(*m) == 0 {
		return h
	}
	for i := len(*m) - 1; i >= 0; i-- {
		h = (*m)[i](h)
	}
	return h
}

func (m *Middlewares) HandleFunc(h http.Handler) http.Handler {
	return m.Handle(h)
}
