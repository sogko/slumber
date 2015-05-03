package server

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"net/http"
)

// RouteHandlerVersion type
type RouteHandlerVersion string

// RouteHandlers is a map of route version to its handler
type RouteHandlers map[RouteHandlerVersion]http.HandlerFunc

// Route type
// Note that DefaultVersion must exists in RouteHandlers map
// See server/routes.go for examples
type Route struct {
	Name           string
	Method         string
	Pattern        string
	DefaultVersion RouteHandlerVersion
	RouteHandlers  RouteHandlers
}

// Routes type
type Routes []Route

// Router type
type Router struct {
	*mux.Router
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

		defaultHandler, ok := route.RouteHandlers[route.DefaultVersion]
		if !ok {
			// server/router instantiation error
			// its safe to throw panic here
			panic(errors.New(fmt.Sprintf("Routes definition error, missing default route handler for version `%v` in `%v`",
				route.DefaultVersion, route.Name)))
		}

		// matcherFunc matches the handler to the correct API version based on its `accept` header
		matcherFunc := func(r *http.Request, rm *mux.RouteMatch) bool {

			headers := utils.ParseAcceptHeaders(r.Header.Get("accept"))
			rm.Handler = defaultHandler

			// try to match a handler to the specified `version` params
			// else we will fall back to the default handler
			for _, h := range headers {
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
				if handler, ok := route.RouteHandlers[RouteHandlerVersion(version)]; ok {
					// found handler for specified version
					rm.Handler = handler
					break
				}
			}
			return true
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			MatcherFunc(matcherFunc)

	}
	return &Router{router}
}
