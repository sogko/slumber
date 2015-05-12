package domain

import (
	"net/http"
)

type ContextHandlerFunc func(http.ResponseWriter, *http.Request, IContext)

func (h ContextHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h(rw, r, nil)
}

type IMiddleware interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
}
