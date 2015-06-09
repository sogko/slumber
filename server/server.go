package server

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"net/http"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	negroni *negroni.Negroni
	Context domain.IContext
	router  *Router
}

// Config type
type Config struct {
	Database        domain.IDatabaseOptions
	Renderer        domain.IRendererOptions
	Routes          *domain.Routes
	TokenAuthority  domain.ITokenAuthorityOptions
	ACLMap          *domain.ACLMap
	ControllerHooks *domain.ControllerHooksMap
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up server and middlewares
	n := negroni.Classic()

	// set up request context
	ctx := middlewares.NewContext()

	s := &Server{n, ctx, nil}

	// set up AccessController
	ac := middlewares.NewAccessController()
	if options.ACLMap != nil {
		ac.Add(options.ACLMap)
	}

	// set up router
	// TODO: de-couple ctx and ac from router
	s.router = NewRouter(options.Routes, ctx, ac)

	// set up db session
	if options.Database != nil {
		db := middlewares.NewMongoDB(options.Database)
		dbSession := db.NewSession()
		s.UseMiddleware(dbSession.Handler)
	}
	// set up renderer
	if options.Renderer != nil {
		renderer := middlewares.NewRenderer(options.Renderer)
		s.UseMiddleware(renderer.Handler)
	}
	// set up TokenAuthority
	if options.TokenAuthority != nil {
		ta := middlewares.NewTokenAuthority(options.TokenAuthority)
		s.UseMiddleware(ta.Handler)
	}
	// set up Authenticator
	auth := middlewares.NewAuthenticator()
	s.UseMiddleware(auth.Handler)

	// register controller hooks
	if options.ControllerHooks == nil {
		options.ControllerHooks = &domain.ControllerHooksMap{}
	}
	hooks := middlewares.ControllerHooksMiddleware{options.ControllerHooks}
	s.UseMiddleware(hooks.Handler)

	return s
}

func (s *Server) UseMiddleware(middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx domain.IContext)) *Server {
	s.negroni.Use(negroni.HandlerFunc(s.Context.InjectWithNext(middleware)))
	return s
}

func (s *Server) SetupRoutes() *Server {
	// add router and clear mux.context values at the end of request life-times
	s.negroni.UseHandler(context.ClearHandler(s.router))
	return s
}

func (s *Server) Run(address string) *Server {
	s.negroni.Run(address)
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) *Server {
	s.negroni.ServeHTTP(w, r)
	return s
}
