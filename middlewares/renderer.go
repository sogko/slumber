package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/unrolled/render"
	"net/http"
)

// RendererOptions type
type RendererOptions render.Options

// Renderer type
type Renderer struct {
	*render.Render
}

// NewRenderer Returns a new Renderer object
func NewRenderer(options domain.IRendererOptions) domain.IRenderer {
	r := render.New(render.Options(*options.(*RendererOptions)))
	return &Renderer{r}
}

// HandlerWithNext Returns a middleware HandlerFunc that saves the Render object into request context
func (renderer *Renderer) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	// create a new renderer and save it in the  request context
	// unrolled/render is a global object that is thread-safe by desi
	ctx.SetRendererCtx(req, renderer)
	next(w, req)
}

func (renderer *Renderer) JSON(w http.ResponseWriter, status int, v interface{}) {
	renderer.Render.JSON(w, status, v)
}

func (renderer *Renderer) Data(w http.ResponseWriter, status int, v []byte) {
	renderer.Render.Data(w, status, v)
}
