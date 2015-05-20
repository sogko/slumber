package domain

import (
	"net/http"
)

type IRendererOptions interface {
}

// Renderer interface
type IRenderer interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
	JSON(w http.ResponseWriter, status int, v interface{})
	Data(w http.ResponseWriter, status int, v []byte)
	Text(w http.ResponseWriter, status int, v []byte)
}
