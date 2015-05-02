package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Route type
// Example of a route:
// Name: 		"CustomersList"
// Method: 		"GET"
// Pattern: 	"/customers
// HandlerFunc: func HandleCustomersGet(w http.ResponseWriter, req *http.Request) { ... }
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes type
type Routes []Route

// Router type
type Router struct {
	*mux.Router
}

// NewRouter Returns a new Router object
func NewRouter(routes *Routes) *Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return &Router{router}
}
