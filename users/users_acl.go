package users

import (
	"github.com/gorilla/mux"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/repositories"
	"net/http"
)

var ACL = domain.ACLMap{
	ListUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		return true, ""
	},
	GetUser: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		return true, ""
	},
	CreateUser: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		// allow anonymous to create a user account
		// if authenticated, only admin can create new users
		// no point for non-admin to create new users
		// TODO: only allow authorized but unauthenticated client
		if user == nil {
			// enforce authenticated access
			return true, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if !user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return false, ""
		}
		return true, ""
	},
	UpdateUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if !user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return false, ""
		}
		// only logged-in admins can update users in batch
		return true, ""
	},
	DeleteAllUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if !user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return false, ""
		}
		// only logged-in admins can update users in batch
		return true, ""
	},
	ConfirmUser: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		// allow anonymous access. user is expected to specify `code`
		return true, ""
	},
	UpdateUser: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		params := mux.Vars(req)
		id := params["id"]
		db := ctx.GetDbCtx(req)
		repo := repositories.UserRepository{db}

		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return true, ""
		}

		// retrieve target user
		userTarget, _ := repo.GetUserById(id)
		if userTarget != nil && user.ID == userTarget.ID {
			// this is his own account
			return true, ""
		}
		// a user can only `update` its own user account or if user is an admin
		return false, ""
	},
	DeleteUser: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		// only an admin can `delete` a user account
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if !user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return false, ""
		}
		// only logged-in admins can update users in batch
		return true, ""
	},
	CountUsers: func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
		if user == nil {
			// enforce authenticated access
			return false, ""
		}
		if user.Status != domain.StatusActive {
			// must be an active user
			return false, ""
		}
		if !user.HasRole(domain.RoleAdmin) {
			// must have an admin role
			return false, ""
		}
		// only logged-in admins can update users in batch
		return true, ""
	},
}
