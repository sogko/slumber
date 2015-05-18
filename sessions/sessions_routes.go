package sessions

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
)

const (
	GetSession    = "GetSession"
	CreateSession = "CreateSession"
	DeleteSession = "DeleteSession"
)

// SessionsAPIRoutes Wire API routes to controllers (http.HandlerFunc)
var SessionsAPIRoutes = domain.Routes{

	//------------- API /sessions ---------//
	domain.Route{
		Name:           GetSession,
		Method:         "GET",
		Pattern:        "/api/sessions",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleGetSession_v0,
		},
	},
	domain.Route{
		Name:           CreateSession,
		Method:         "POST",
		Pattern:        "/api/sessions",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleCreateSession_v0,
		},
	},
	domain.Route{
		Name:           DeleteSession,
		Method:         "DELETE",
		Pattern:        "/api/sessions",
		DefaultVersion: "0.0",
		RouteHandlers: domain.RouteHandlers{
			"0.0": HandleDeleteSession_v0,
		},
	},
}
