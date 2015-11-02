package server

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/sogko/slumber/domain"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"time"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	negroni        *negroni.Negroni
	Context        domain.IContext
	router         *Router
	gracefulServer *graceful.Server
	timeout        time.Duration
}

// Config type
type Config struct {
	Context domain.IContext
}

// Options for running the server
type Options struct {
	Timeout time.Duration
	ShutdownHandler func()
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up server and middlewares
	n := negroni.Classic()

	s := &Server{n, options.Context, nil, nil, 0}

	return s
}

func (s *Server) UseMiddleware(middleware domain.IMiddleware) *Server {
	// next convert it into negroni style handlerfunc
	s.negroni.Use(negroni.HandlerFunc(middleware.Handler))
	return s
}

func (s *Server) UseContextMiddleware(middleware domain.IContextMiddleware) *Server {
	// take contextual middleware, inject context into it.
	// next convert it into negroni style handlerfunc
	s.negroni.Use(negroni.HandlerFunc(s.Context.InjectMiddleware(middleware.Handler)))
	return s
}

func (s *Server) UseRouter(router *Router) *Server {
	// add router and clear mux.context values at the end of request life-times
	s.negroni.UseHandler(context.ClearHandler(router))
	return s
}

func (s *Server) Run(address string, options Options) *Server {
	s.timeout = options.Timeout
	s.gracefulServer = &graceful.Server{
		Timeout:           options.Timeout,
		Server:            &http.Server{Addr: address, Handler: s.negroni},
		ShutdownInitiated: options.ShutdownHandler,
	}
	s.gracefulServer.ListenAndServe()
	return s
}

func (s *Server) Stop() {
	s.gracefulServer.Stop(s.timeout)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) *Server {
	s.negroni.ServeHTTP(w, r)
	return s
}
