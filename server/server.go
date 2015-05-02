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

// Components type
type Components struct {
	DatabaseSession *DatabaseSession
	Renderer        *Renderer
}

// NewServer Returns a new Server object
func NewServer(components *Components) *Server {

	// set up router
	routes := GetRoutes()
	r := NewRouter(routes)

	// set up server and middlewares
	n := negroni.Classic()
	n.Use(components.DatabaseSession.UseDatabase())
	n.Use(components.Renderer.UseRenderer())

	// add router and clear mux.context values at the end of request life-times
	n.UseHandler(context.ClearHandler(r))

	return &Server{n}
}
