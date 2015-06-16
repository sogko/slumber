package domain

import "net/http"

// RouteHandlerVersion type
type RouteHandlerVersion string

// RouteHandlers is a map of route version to its handler
type RouteHandlers map[RouteHandlerVersion]http.HandlerFunc

// Route type
// Note that DefaultVersion must exists in RouteHandlers map
// See routes.go for examples
type Route struct {
	Name           string
	Method         string
	Pattern        string
	DefaultVersion RouteHandlerVersion
	RouteHandlers  RouteHandlers
	ACLHandler     ACLHandlerFunc
}

// Routes type
type Routes []Route

// Append Returns a new slice of Routes
func (r *Routes) Append(routes ...*Routes) Routes {
	res := Routes{}
	// copy current route
	for _, route := range *r {
		res = append(res, route)
	}
	for _, _routes := range routes {
		for _, route := range *_routes {
			res = append(res, route)
		}
	}
	return res
}
