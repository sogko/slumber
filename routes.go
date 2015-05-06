package main

import (
	"github.com/sogko/golang-rest-api-server-example/server"
	"github.com/sogko/golang-rest-api-server-example/sessions"
	"github.com/sogko/golang-rest-api-server-example/users"
)

// GetRoutes Wire API routes to controllers (http.HandlerFunc)
func GetRoutes() *server.Routes {

	return &server.Routes{
		//------------- API /users ---------//
		server.Route{
			Name:           "ListUsers",
			Method:         "GET",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleList_v0,
			},
		},
		server.Route{
			Name:           "CreateUser",
			Method:         "POST",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleCreate_v0,
			},
		},
		server.Route{
			Name:           "UpdateUsers",
			Method:         "PUT",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleUpdateList_v0,
			},
		},
		server.Route{
			Name:           "DeleteAllUsers",
			Method:         "DELETE",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleDeleteAll_v0,
			},
		},
		server.Route{
			Name:           "GetUser",
			Method:         "GET",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleGet_v0,
			},
		},
		/*
			Method for email confirmation has to be GET because
			link to confirm email has to be click-able from email content
			(You can't add a POST/PUT body)
		*/
		server.Route{
			Name:           "ConfirmUser",
			Method:         "GET",
			Pattern:        "/api/users/{id}/confirm",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleConfirmUser_v0,
			},
		},
		server.Route{
			Name:           "UpdateUser",
			Method:         "PUT",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleUpdate_v0,
			},
		},
		server.Route{
			Name:           "DeleteUser",
			Method:         "DELETE",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": users.HandleDelete_v0,
			},
		},
		//------------- API /users ---------//
		server.Route{
			Name:           "CreateSession",
			Method:         "POST",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: server.RouteHandlers{
				"0.0": sessions.HandleCreate_v0,
			},
		},
	}
}
