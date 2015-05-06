//TODO: add ACL to Users API
// Only admin are able to access Users API
package users

import (
	"github.com/gorilla/mux"
	"github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
)

// HandleList_v0 lists users
func HandleList_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)

	users := GetUsers(db)

	r.JSON(w, http.StatusOK, ListResponse_v0{
		Users:   users,
		Message: "User list retrieved",
		Success: true,
	})
}

// HandleUpdateList_v0 update a list of users
func HandleUpdateList_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)

	var body UpdateListRequest_v0
	err := server.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	var message = "User list updated"
	var success bool = true

	if body.Action == "delete" {
		err = DeleteUsers(db, body.IDs)
	} else {
		r.JSON(w, http.StatusBadRequest, UpdateListResponse_v0{
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

	r.JSON(w, http.StatusOK, UpdateListResponse_v0{
		Action:  body.Action,
		IDs:     body.IDs,
		Message: message,
		Success: success,
	})
}

// HandleDeleteAll_v0 deletes all users
func HandleDeleteAll_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)

	_ = DeleteAllUsers(db)

	r.JSON(w, http.StatusOK, DeleteAllResponse_v0{
		Message: "All users deleted",
		Success: true,
	})
}

// HandleCreate_v0 creates a new user
func HandleCreate_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)

	var body CreateRequest_v0
	err := server.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	// New user always have no roles assigned until confirmed
	body.User.Roles = Roles{}

	// Set flag to `pending` awaiting user to confirm email
	body.User.Status = StatusPending

	// generate new code
	body.User.GenerateConfirmationCode()

	// ensure that user obj is valid
	if !body.User.IsValid() {
		server.RenderErrorResponseHelper(w, req, r, "Invalid user object")
		return
	}

	err = CreateUser(db, &body.User)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, "Failed to save user object")
		return
	}

	// TODO: send email / message with eamil confirmation code

	r.JSON(w, http.StatusCreated, CreateResponse_v0{
		User:    body.User,
		Message: "User created",
		Success: true,
	})
}

// HandleConfirmEmail_v0 confirms user's email address
func HandleConfirmUser_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)
	params := mux.Vars(req)
	id := params["id"]
	code := req.FormValue("code")

	user, err := GetUser(db, id)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	if user.Status != StatusPending {
		r.JSON(w, http.StatusBadRequest, ConfirmUserResponse_v0{
			Code:    code,
			User:    *user,
			Message: "User not pending confirmation",
			Success: false,
		})
		return
	}

	if !user.IsCodeVerified(code) {
		r.JSON(w, http.StatusBadRequest, ConfirmUserResponse_v0{
			Code:    code,
			User:    *user,
			Message: "Invalid code",
			Success: false,
		})
		return
	}

	// set user status to `active`
	user, err = UpdateUser(db, id, &User{
		Status: StatusActive,
		Roles:  Roles{RoleUser},
	})
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, ConfirmUserResponse_v0{
		Code:    code,
		User:    *user,
		Message: "User confirmed",
		Success: true,
	})
}

// HandleGet_v0 gets user object
func HandleGet_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	user, err := GetUser(db, id)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, "User not found")
		return
	}

	r.JSON(w, http.StatusOK, GetResponse_v0{
		User:    *user,
		Message: "User retrieved",
		Success: true,
	})
}

// HandleUpdate_v0 updates user object
func HandleUpdate_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	var body UpdateRequest_v0
	err := server.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	user, err := UpdateUser(db, id, &body.User)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, UpdateResponse_v0{
		User:    *user,
		Message: "User updated",
		Success: true,
	})
}

// HandleDelete_v0 deletes user object
func HandleDelete_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	err := DeleteUser(db, id)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	r.JSON(w, http.StatusOK, DeleteResponse_v0{
		Message: "User deleted",
		Success: true,
	})
}
