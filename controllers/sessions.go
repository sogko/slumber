package controllers

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"log"
	"net/http"
)

type CreateSessionRequest_v0 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type CreateSessionResponse_v0 struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
}

// HandleCreateSession_v0 verify user's credentials and generates a JWT token if valid
func HandleCreateSession_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	ta := ctx.GetTokenAuthorityCtx(req)

	var body CreateSessionRequest_v0
	err := DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	if body.Username == "" {
		RenderErrorResponseHelper(w, req, r, "Empty username")
		return
	}

	repo := repositories.UserRepository{db}
	user, err := repo.GetUserByUsername(body.Username)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	if !user.IsCredentialsVerified(body.Password) {
		RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	var rolesString []string
	for _, role := range user.Roles {
		rolesString = append(rolesString, string(role))
	}

	tokenString, err := ta.CreateNewSessionToken(&domain.TokenClaims{
		UserID: user.ID.Hex(),
		Status: user.Status,
		Roles:  rolesString,
	})

	if err != nil {
		log.Println("err", err.Error())
		RenderErrorResponseHelper(w, req, r, "Error creating session token")
		return
	}

	r.JSON(w, http.StatusCreated, CreateSessionResponse_v0{
		Token:   tokenString,
		Success: true,
	})
}

// HandleDeleteSession_v0 invalidates a session token
func HandleDeleteSession_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	claims := ctx.GetAuthenticatedClaimsCtx(req)

	log.Println("Claim", claims)
	var body CreateSessionRequest_v0
	err := DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	r.JSON(w, http.StatusOK, CreateSessionResponse_v0{
		Token:   "TEST",
		Success: true,
	})
}
