package server

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/middlewares"
)

// Request JSON body limit is set at 5MB (currently not enforced)
const BodyLimitBytes uint32 = 1048576 * 5

// Server type
type Server struct {
	*negroni.Negroni
}

// Config type
type Config struct {
	Database       domain.IDatabaseOptions
	Renderer       domain.IRendererOptions
	Routes         *domain.Routes
	TokenAuthority domain.ITokenAuthorityOptions
	ACLMap         *domain.ACLMap
}

// NewServer Returns a new Server object
func NewServer(options *Config) *Server {

	// set up request context
	ctx := middlewares.NewContext()

	// set up AccessController
	ac := middlewares.NewAccessController()
	if options.ACLMap != nil {
		ac.Add(options.ACLMap)
	}

	// set up router
	// TODO: de-couple ctx and ac from router
	router := NewRouter(options.Routes, ctx, ac)

	// set up db session
	db := middlewares.NewMongoDB(options.Database)
	dbSession := db.NewSession()

	// set up renderer
	renderer := middlewares.NewRenderer(options.Renderer)

	// set up TokenAuthority
	ta := middlewares.NewTokenAuthority(options.TokenAuthority)

	// set up Authenticator
	auth := middlewares.NewAuthenticator()

	// set up server and middlewares
	n := negroni.Classic()
	n.Use(negroni.HandlerFunc(ctx.InjectWithNext(dbSession.Handler)))
	n.Use(negroni.HandlerFunc(ctx.InjectWithNext(renderer.Handler)))
	n.Use(negroni.HandlerFunc(ctx.InjectWithNext(ta.Handler)))
	n.Use(negroni.HandlerFunc(ctx.InjectWithNext(auth.Handler)))

	// add router and clear mux.context values at the end of request life-times
	n.UseHandler(context.ClearHandler(router))

	return &Server{n}
}
