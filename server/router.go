package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/libs"
	"net/http"
)

// Router type
type Router struct {
	*mux.Router
}

// matcherFunc matches the handler to the correct API version based on its `accept` header
// TODO: refactor matcher function as server.Config
func matcherFunc(r domain.Route, defaultHandler domain.ContextHandlerFunc, ctx domain.IContext, ac domain.IAccessController) func(r *http.Request, rm *mux.RouteMatch) bool {
	return func(req *http.Request, rm *mux.RouteMatch) bool {
		acceptHeaders := libs.ParseAcceptHeaders(req.Header.Get("accept"))
		foundHandler := defaultHandler
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
			if handler, ok := r.RouteHandlers[domain.RouteHandlerVersion(version)]; ok {
				// found handler for specified version
				foundHandler = handler
				break
			}
		}
		// injects ctx IContext and add AccessController into routeHandler
		rm.Handler = ctx.Inject(ac.Handler(r.Name, foundHandler))
		return true
	}
}

// NewRouter Returns a new Router object
func NewRouter(routes *domain.Routes, ctx domain.IContext, ac domain.IAccessController) *Router {
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
			MatcherFunc(matcherFunc(route, defaultHandler, ctx, ac))

	}
	return &Router{router}
}
