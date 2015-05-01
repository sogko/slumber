package server

import (
	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	"net/http"
)

type Render render.Render

type Renderer struct {
	Render *render.Render
}

//
func NewRenderer(options render.Options) *Renderer {
	r := render.New(options)
	return &Renderer{r}
}

// Returns a negroni middleware HandlerFunc that saves the Render object into request context
func (renderer *Renderer) UseRenderer() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// create a new renderer and save it in the  request context
		SetRenderCtx(r, renderer.Render)
		next(rw, r)
	})
}
