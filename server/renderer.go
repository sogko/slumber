package server

import (
	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	"net/http"
)

// Renderer type
type Renderer struct {
	*render.Render
}

// NewRenderer Returns a new Renderer object
func NewRenderer(options render.Options) *Renderer {
	r := render.New(options)
	return &Renderer{r}
}

// UseRenderer Returns a negroni middleware HandlerFunc that saves the Render object into request context
func (renderer *Renderer) UseRenderer() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// create a new renderer and save it in the  request context

		SetRendererCtx(r, renderer)
		next(rw, r)
	})
}
