package users

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
)

const (
	ListUsers      = "ListUsers"
	CountUsers     = "CountUsers"
	GetUser        = "GetUser"
	CreateUser     = "CreateUser"
	UpdateUsers    = "UpdateUsers"
	DeleteAllUsers = "DeleteAllUsers"
	ConfirmUser    = "ConfirmUser"
	UpdateUser     = "UpdateUser"
	DeleteUser     = "DeleteUser"
)

// UsersAPIRoutes Wire API routes to controllers (http.HandlerFunc)
var UsersAPIRoutes = domain.Routes{
	domain.Route{
		Name:           ListUsers,
		Method:         "GET",
		Pattern:        "/api/users",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleListUsers_v0,
		},
	},
	domain.Route{
		Name:           CountUsers,
		Method:         "GET",
		Pattern:        "/api/users/count",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleCountUsers_v0,
		},
	},
	domain.Route{
		Name:           CreateUser,
		Method:         "POST",
		Pattern:        "/api/users",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleCreateUser_v0,
		},
	},
	domain.Route{
		Name:           UpdateUsers,
		Method:         "PUT",
		Pattern:        "/api/users",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleUpdateUsers_v0,
		},
	},
	domain.Route{
		Name:           DeleteAllUsers,
		Method:         "DELETE",
		Pattern:        "/api/users",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleDeleteAllUsers_v0,
		},
	},
	domain.Route{
		Name:           GetUser,
		Method:         "GET",
		Pattern:        "/api/users/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleGetUser_v0,
		},
	},
	/*
		Method for email confirmation has to be GET because
		link to confirm email has to be click-able from email content
		(You can't add a POST/PUT body)
	*/
	domain.Route{
		Name:           ConfirmUser,
		Method:         "GET",
		Pattern:        "/api/users/{id}/confirm",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleConfirmUser_v0,
		},
	},
	domain.Route{
		Name:           UpdateUser,
		Method:         "PUT",
		Pattern:        "/api/users/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleUpdateUser_v0,
		},
	},
	domain.Route{
		Name:           DeleteUser,
		Method:         "DELETE",
		Pattern:        "/api/users/{id}",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleDeleteUser_v0,
		},
	},
}
