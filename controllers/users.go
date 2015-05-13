package controllers

import (
	"github.com/gorilla/mux"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"net/http"
)

//---- User Request API v0 ----

type ListUsersResponse_v0 struct {
	Users   domain.Users `json:"users"`
	Message string       `json:"message,omitempty"`
	Success bool         `json:"success"`
}

type CreateUserRequest_v0 struct {
	User domain.NewUser `json:"user"`
}

type CreateUserResponse_v0 struct {
	User    domain.User `json:"user,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

type ConfirmUserResponse_v0 struct {
	Code    string      `json:"code,omitempty"`
	User    domain.User `json:"user,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}

type UpdateUsersRequest_v0 struct {
	Action string   `json:"action"`
	IDs    []string `json:"ids"`
}

type UpdateUsersResponse_v0 struct {
	Action  string   `json:"action,omitempty"`
	IDs     []string `json:"ids,omitempty"`
	Message string   `json:"message,omitempty"`
	Success bool     `json:"success"`
}

type DeleteAllUsersResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}
type GetUserResponse_v0 struct {
	User    domain.User `json:"user,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}
type UpdateUserRequest_v0 struct {
	User domain.User `json:"user"`
}

type UpdateUserResponse_v0 struct {
	User    domain.User `json:"user,omitempty"`
	Message string      `json:"message,omitempty"`
	Success bool        `json:"success"`
}
type DeleteUserResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// HandleListUsers_v0 lists users
func HandleListUsers_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)

	repo := repositories.UserRepository{db}
	users := repo.GetUsers()

	r.JSON(w, http.StatusOK, ListUsersResponse_v0{
		Users:   users,
		Message: "User list retrieved",
		Success: true,
	})
}

// HandleUpdateList_v0 update a list of users
func HandleUpdateUsers_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)

	var body UpdateUsersRequest_v0
	err := DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	var message = "User list updated"
	var success bool = true

	if body.Action == "delete" {
		repo := repositories.UserRepository{db}
		err = repo.DeleteUsers(body.IDs)
	} else {
		r.JSON(w, http.StatusBadRequest, UpdateUsersResponse_v0{
			Action:  body.Action,
			IDs:     body.IDs,
			Message: "Invalid action",
			Success: false,
		})
		return
	}
	if err != nil {
		success = false
		message = err.Error()
	}

	r.JSON(w, http.StatusOK, UpdateUsersResponse_v0{
		Action:  body.Action,
		IDs:     body.IDs,
		Message: message,
		Success: success,
	})
}

// HandleDeleteAll_v0 deletes all users
func HandleDeleteAllUsers_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)

	repo := repositories.UserRepository{db}
	_ = repo.DeleteAllUsers()

	r.JSON(w, http.StatusOK, DeleteAllUsersResponse_v0{
		Message: "All users deleted",
		Success: true,
	})
}

// HandleCreateUser_v0 creates a new user
func HandleCreateUser_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	repo := repositories.UserRepository{db}

	var body CreateUserRequest_v0
	err := DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	if repo.UserExistsByUsername(body.User.Username) {
		RenderErrorResponseHelper(w, req, r, "Username already exists")
		return
	}

	if repo.UserExistsByEmail(body.User.Email) {
		RenderErrorResponseHelper(w, req, r, "User with email address already exists")
		return
	}

	// New user always have no roles assigned until confirmed
	// Set flag to `pending` awaiting user to confirm email
	var newUser = domain.User{
		Username: body.User.Username,
		Email:    body.User.Email,
		Roles:    domain.Roles{},
		Status:   domain.StatusPending,
	}

	// generate new code
	newUser.GenerateConfirmationCode()

	// set password (hashed)
	newUser.SetPassword(body.User.Password)

	// ensure that user obj is valid
	if !newUser.IsValid() {
		RenderErrorResponseHelper(w, req, r, "Invalid user object")
		return
	}

	err = repo.CreateUser(&newUser)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, "Failed to save user object")
		return
	}

	// TODO: send email / message with email confirmation code

	r.JSON(w, http.StatusCreated, CreateUserResponse_v0{
		User:    newUser,
		Message: "User created",
		Success: true,
	})
}

// HandleConfirmEmail_v0 confirms user's email address
func HandleConfirmUser_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]
	code := req.FormValue("code")

	repo := repositories.UserRepository{db}
	user, err := repo.GetUserById(id)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	if user.Status != domain.StatusPending {
		RenderErrorResponseHelper(w, req, r, "User not pending confirmation")
		return
	}

	if !user.IsCodeVerified(code) {
		RenderErrorResponseHelper(w, req, r, "Invalid code")
		return
	}

	// set user status to `active`
	user, err = repo.UpdateUser(id, &domain.User{
		Status: domain.StatusActive,
		Roles:  domain.Roles{domain.RoleUser},
	})
	if err != nil {
		RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, ConfirmUserResponse_v0{
		Code:    code,
		User:    *user,
		Message: "User confirmed",
		Success: true,
	})
}

// HandleGetUser_v0 gets user object
func HandleGetUser_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	repo := repositories.UserRepository{db}
	user, err := repo.GetUserById(id)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, "User not found")
		return
	}

	r.JSON(w, http.StatusOK, GetUserResponse_v0{
		User:    *user,
		Message: "User retrieved",
		Success: true,
	})
}

// HandleUpdateUser_v0 updates user object
func HandleUpdateUser_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	var body UpdateUserRequest_v0
	err := DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	repo := repositories.UserRepository{db}
	user, err := repo.UpdateUser(id, &body.User)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, UpdateUserResponse_v0{
		User:    *user,
		Message: "User updated",
		Success: true,
	})
}

// HandleDelete_v0 deletes user object
func HandleDeleteUser_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	repo := repositories.UserRepository{db}

	err := repo.DeleteUser(id)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, DeleteUserResponse_v0{
		Message: "User deleted",
		Success: true,
	})
}
