package sessions

import (
	"github.com/sogko/slumber/domain"
	"net/http"
)

var SessionsAPIACL = domain.ACLMap{
	GetSession: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		if user == nil {
			return false
		}
		return true
	},
	CreateSession: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// allow anonymous access
		return true
	},
	DeleteSession: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// allow anonymous access
		return true
	},
}
