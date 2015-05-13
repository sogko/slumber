package acl

import (
	"github.com/gorilla/mux"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"net/http"
)

const (
	ListUsers      = "ListUsers"
	GetUser        = "GetUser"
	CreateUser     = "CreateUser"
	UpdateUsers    = "UpdateUsers"
	DeleteAllUsers = "DeleteAllUsers"
	ConfirmUser    = "ConfirmUser"
	UpdateUser     = "UpdateUser"
	DeleteUser     = "DeleteUser"
)

var UsersAPIACL = domain.ACLMap{
	ListUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// enforce authenticated access
		return (user != nil &&
			user.Status == domain.StatusActive)
	},
	GetUser: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// enforce authenticated access
		return (user != nil &&
			user.Status == domain.StatusActive)
	},
	CreateUser: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// allow anonymous access; anyone can create a user account
		// TODO: only allow authorized but unauthenticated client
		return true
	},
	UpdateUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// only logged-in admins can update users in batch
		return (user != nil &&
			user.Status == domain.StatusActive &&
			user.HasRole(domain.RoleAdmin))
	},
	DeleteAllUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// only logged-in admins can update users in batch
		return (user != nil &&
			user.Status == domain.StatusActive &&
			user.HasRole(domain.RoleAdmin))
	},
	ConfirmUser: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// allow anonymous access. user is expected to specify `code`
		return true
	},
	UpdateUser: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {

		// enforce authenticated access
		if user == nil {
			return false
		}

		params := mux.Vars(req)
		id := params["id"]
		db := ctx.GetDbCtx(req)
		repo := repositories.UserRepository{db}

		// retrieve target user
		userTarget, err := repo.GetUserById(id)
		if err != nil {
			return false
		}

		// a user can only `update` its own user account or if user is an admin
		return (user.Status == domain.StatusActive &&
			(user.ID == userTarget.ID || user.HasRole(domain.RoleAdmin)))
	},
	DeleteUser: func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
		// only an admin can `delete` a user account
		return (user != nil &&
			user.Status == domain.StatusActive &&
			user.HasRole(domain.RoleAdmin))
	},
}
