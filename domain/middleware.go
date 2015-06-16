package domain

import (
	"net/http"
)

type ContextHandlerFunc func(http.ResponseWriter, *http.Request, IContext)

func (h ContextHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, ctx IContext) {
	h(rw, r, ctx)
}

type MiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (m MiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	m(rw, r, next)
}

type ContextMiddlewareFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)

func (m ContextMiddlewareFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext) {
	m(rw, r, next, ctx)
}

type IMiddleware interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type IContextMiddleware interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
}
