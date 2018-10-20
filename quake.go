package quake

import (
	"net/http"
)

//MiddlewareFunc ...
type MiddlewareFunc func(http.Handler) http.Handler

//Middleware ...
type Middleware struct {
	stack []MiddlewareFunc
	final http.Handler
}

//NewMiddleware ...
func NewMiddleware() *Middleware {
	return &Middleware{
		stack: make([]MiddlewareFunc, 0),
	}
}

//Add a new middleware function to the stack
func (mw *Middleware) Add(mdfunc ...MiddlewareFunc) {
	for _, h := range mdfunc {
		mw.stack = append(mw.stack, h)
	}
}

//ServeHTTP implements http.Handler
func (mw *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.final.ServeHTTP(w, r)
}

//Wrap receives an http.Handler and applies the middleware chain to it
func (mw *Middleware) Wrap(final http.Handler) {
	mw.final = final
	for i := len(mw.stack); i > 0; i-- {
		mw.final = mw.stack[i-1](mw.final)
	}
}
