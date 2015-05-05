package server

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	*negroni.Negroni
}

// Config type
type Config struct {
	Database *DatabaseOptions
	Renderer *RendererOptions
	Routes   *Routes
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up router
	router := NewRouter(options.Routes)

	// set up db session
	session := NewSession(*options.Database)

	// set up renderer
	renderer := NewRenderer(*options.Renderer)

	// set up server and middlewares
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(session.HandlerWithNext))
	n.Use(negroni.HandlerFunc(renderer.HandlerWithNext))

	// add router and clear mux.context values at the end of request life-times
	n.UseHandler(context.ClearHandler(router))

	return &Server{n}
}
