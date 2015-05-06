package sessions

import (
	"github.com/sogko/golang-rest-api-server-example/server"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"net/http"
)

type CreateRequest_v0 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type CreateResponse_v0 struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}


// HandleCreate_v0 verify user's credentials and generates a JWT token if valid
func HandleCreate_v0(w http.ResponseWriter, req *http.Request) {
	r := server.RendererCtx(req)
	db := server.DbCtx(req)

	var body CreateRequest_v0
	err := server.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	if body.Username == "" {
		server.RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	user, err := GetUserByUsername(db, body.Username)
	if err != nil {
		server.RenderErrorResponseHelper(w, req, r, err.Error())
		return
	}

	if !user.IsCredentialsVerified(body.Password) {
		server.RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	var rolesString []string
	for _, role := range user.Roles {
		rolesString = append(rolesString, string(role))
	}

	tokenString, _ := utils.CreateNewToken(utils.TokenClaims{
		ID:     user.ID.Hex(),
		Status: user.Status,
		Roles:  rolesString,
	})

	r.JSON(w, http.StatusCreated, CreateResponse_v0{
		Token:   tokenString,
		Success: true,
	})
}
