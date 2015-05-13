package main

import (
	"github.com/sogko/golang-rest-api-server-example/acl"
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
)

// GetRoutes Wire API routes to controllers (http.HandlerFunc)
func GetRoutes() *domain.Routes {

	return &domain.Routes{
		domain.Route{
			Name:           "ListUsers",
			Method:         "GET",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleListUsers_v0,
			},
			ACLAction: acl.ListUsers,
		},
		domain.Route{
			Name:           "CreateUser",
			Method:         "POST",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleCreateUser_v0,
			},
			ACLAction: acl.CreateUser,
		},
		domain.Route{
			Name:           "UpdateUsers",
			Method:         "PUT",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleUpdateUsers_v0,
			},
			ACLAction: acl.UpdateUsers,
		},
		domain.Route{
			Name:           "DeleteAllUsers",
			Method:         "DELETE",
			Pattern:        "/api/users",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleDeleteAllUsers_v0,
			},
			ACLAction: acl.DeleteAllUsers,
		},
		domain.Route{
			Name:           "GetUser",
			Method:         "GET",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleGetUser_v0,
			},
			ACLAction: acl.GetUser,
		},
		/*
			Method for email confirmation has to be GET because
			link to confirm email has to be click-able from email content
			(You can't add a POST/PUT body)
		*/
		domain.Route{
			Name:           "ConfirmUser",
			Method:         "GET",
			Pattern:        "/api/users/{id}/confirm",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleConfirmUser_v0,
			},
			ACLAction: acl.ConfirmUser,
		},
		domain.Route{
			Name:           "UpdateUser",
			Method:         "PUT",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleUpdateUser_v0,
			},
			ACLAction: acl.UpdateUser,
		},
		domain.Route{
			Name:           "DeleteUser",
			Method:         "DELETE",
			Pattern:        "/api/users/{id}",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleDeleteUser_v0,
			},
			ACLAction: acl.DeleteUser,
		},
		//------------- API /sessions ---------//
		domain.Route{
			Name:           "CreateSession",
			Method:         "POST",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleCreateSession_v0,
			},
			ACLAction: acl.CreateSession,
		},
		domain.Route{
			Name:           "DeleteSession",
			Method:         "DELETE",
			Pattern:        "/api/sessions",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": controllers.HandleDeleteSession_v0,
			},
			ACLAction: acl.DeleteSession,
		},
	}
}
