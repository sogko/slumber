package domain

import (
	"net/http"
)

type IAuthenticator interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
}
