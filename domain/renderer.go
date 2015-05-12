package domain

import (
	"net/http"
)

type IRendererOptions interface {
}

// Renderer interface
type IRenderer interface {
	NewRenderer(IRendererOptions) IRenderer
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
	JSON(w http.ResponseWriter, status int, v interface{})
}
