package sessions

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
	"net/http"
)

var SessionsAPIACL = domain.ACLMap{
	CreateSession: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// allow anonymous access
		return true
	},
	DeleteSession: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// enforce authenticated access
		return (user != nil)
	},
}
