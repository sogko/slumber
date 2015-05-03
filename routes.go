package main

import (
	"github.com/sogko/golang-rest-api-server-example/customers"
	"github.com/sogko/golang-rest-api-server-example/server"
)

// GetRoutes Wire API routes to controllers (http.HandlerFunc)
func GetRoutes() *server.Routes {
	return &server.Routes{
		server.Route{"CustomersList", "GET", "/api/customers", "0.1", server.RouteHandlers{
			"0.1": customers.HandleCustomersGet,
		}},
		server.Route{"CustomerCreate", "POST", "/api/customers", "0.1", server.RouteHandlers{
			"0.1": customers.HandleCustomersPost,
		}},
		server.Route{"CustomerGet", "GET", "/api/customers/{id}", "0.1", server.RouteHandlers{
			"0.1": customers.HandleCustomerGet,
		}},
	}
}
