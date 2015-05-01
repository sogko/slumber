package server

import (
	"github.com/codegangsta/negroni"
)

const BODY_LIMIT_BYTES = 1048576 * 5 // 5 MB

type Server struct {
	*negroni.Negroni
}

type ServerComponents struct {
	DatabaseSession *DatabaseSession
	Renderer        *Renderer
}

func NewServer(components *ServerComponents) *Server {

	// set up router
	routes := LoadRoutes()
	r := NewRouter(routes)

	// set up server and middlewares
	n := negroni.Classic()
	n.Use(components.DatabaseSession.UseDatabase())
	n.Use(components.Renderer.UseRenderer())

	// add router
	// Note: context.ClearHandler(r) automatically called by gorilla/mux
	n.UseHandler(r)

	return &Server{n}
}
