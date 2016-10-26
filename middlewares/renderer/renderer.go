package renderer

import (
	"github.com/sogko/slumber/domain"
	"github.com/unrolled/render"
	"net/http"
)

const RendererKey domain.ContextKey = "slumber-mddlwr-unrolled-render-key"
const JSON = "json"
const XML = "xml"
const Data = "octet-stream"
const Text = "text"

type Options render.Options

// Renderer type
// implements IRenderer and IContextMiddleware
type Renderer struct {
	r                 *render.Render
	options           *Options
	DefaultRenderType string
}

// New( Returns a new Renderer object
func New(options *Options, defaultRenderType string) *Renderer {
	r := render.New(render.Options(*options))
	return &Renderer{r, options, defaultRenderType}
}

// HandlerWithNext Returns a middleware HandlerFunc that saves the Render object into request context
func (renderer *Renderer) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	SetRendererCtx(ctx, req, renderer)
	next(w, req)
}

func (renderer *Renderer) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	acceptHeaders := domain.NewAcceptHeadersFromString(req.Header.Get("accept"))

	renderType := renderer.DefaultRenderType
	for _, h := range acceptHeaders {
		m := h.MediaType
		if m.SubType == JSON || m.Suffix == JSON {
			renderType = JSON
			break
		}
		if m.SubType == XML || m.Suffix == XML {
			renderType = XML
			break
		}
		if m.SubType == Data || m.Suffix == Data {
			renderType = Data
			break
		}
		if m.SubType == Text || m.Suffix == Text {
			renderType = Text
			break
		}
	}
	switch renderType {
	case JSON:
		renderer.JSON(w, status, v)
	case XML:
		renderer.XML(w, status, v)
	case Data:
		renderer.Data(w, status, v.([]byte))
	case Text:
		renderer.Text(w, status, v.([]byte))
	default:
		renderer.Text(w, status, v.([]byte))
	}
}

func (renderer *Renderer) JSON(w http.ResponseWriter, status int, v interface{}) {
	renderer.r.JSON(w, status, v)
}

func (renderer *Renderer) XML(w http.ResponseWriter, status int, v interface{}) {
	renderer.r.XML(w, status, v)
}

func (renderer *Renderer) Data(w http.ResponseWriter, status int, v []byte) {
	renderer.r.Data(w, status, v)
}
func (renderer *Renderer) Text(w http.ResponseWriter, status int, v []byte) {
	w.WriteHeader(status)
	w.Write(v)
}

func SetRendererCtx(ctx domain.IContext, r *http.Request, renderer *Renderer) {
	ctx.Set(r, RendererKey, renderer)
}

func GetRendererCtx(ctx domain.IContext, r *http.Request) *Renderer {
	if ren := ctx.Get(r, RendererKey); ren != nil {
		return ren.(*Renderer)
	}
	return nil
}
