package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/middlewares"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"net/http"
)

// RouteHandlerVersion type
type RouteHandlerVersion string

// RouteHandlers is a map of route version to its handler
type RouteHandlers map[RouteHandlerVersion]domain.ContextHandlerFunc

// Route type
// Note that DefaultVersion must exists in RouteHandlers map
// See routes.go for examples
type Route struct {
	Name                 string
	Method               string
	Pattern              string
	DefaultVersion       RouteHandlerVersion
	RouteHandlers        RouteHandlers
	AccessControlHandler RouteHandlers
}

// Routes type
type Routes []Route

// Router type
type Router struct {
	*mux.Router
}

// matcherFunc matches the handler to the correct API version based on its `accept` header
func matcherFunc(r Route, defaultHandler domain.ContextHandlerFunc) func(r *http.Request, rm *mux.RouteMatch) bool {
	ctx := middlewares.Context{}
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		acceptHeaders := utils.ParseAcceptHeaders(req.Header.Get("accept"))
		rm.Handler = ctx.Inject(defaultHandler)

		// try to match a handler to the specified `version` params
		// else we will fall back to the default handler
		for _, h := range acceptHeaders {
			m := h.MediaType
			// check if media type is `application/json` type or `application/[*]+json` suffix
			if !(m.Type == "application" && (m.SubType == "json" || m.Suffix == "json")) {
				continue
			}

			// if its the right application type, check if a version specified
			version, hasVersion := m.Parameters["version"]
			if !hasVersion {
				continue
			}
			if handler, ok := r.RouteHandlers[RouteHandlerVersion(version)]; ok {
				// found handler for specified version
				rm.Handler = ctx.Inject(handler)
				break
			}
		}
		return true
	}
}

// NewRouter Returns a new Router object
func NewRouter(routes *Routes) *Router {
	if routes == nil {
		// server/router instantiation error
		// its safe to throw panic here
		panic(errors.New(fmt.Sprintf("Routes definition missing")))
	}
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {

		// get the defaultHandler for current route at init time so that we can safely panic
		// if it was not defined
		defaultHandler, ok := route.RouteHandlers[route.DefaultVersion]
		if !ok {
			// server/router instantiation error
			// its safe to throw panic here
			panic(errors.New(fmt.Sprintf("Routes definition error, missing default route handler for version `%v` in `%v`",
				route.DefaultVersion, route.Name)))
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			MatcherFunc(matcherFunc(route, defaultHandler))

	}
	return &Router{router}
}
