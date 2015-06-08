package sessions

import (
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type GetSessionResponse_v0 struct {
	User    domain.User `json:"user"`
	Success bool        `json:"success"`
	Message string      `json:"message"`
}
type CreateSessionRequest_v0 struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type CreateSessionResponse_v0 struct {
	Token   string `json:"token"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}
type DeleteSessionResponse_v0 struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandleGetSession_v0 Get session details
func HandleGetSession_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	user := ctx.GetCurrentUserCtx(req)

	r.JSON(w, http.StatusOK, GetSessionResponse_v0{
		User:    *user,
		Success: true,
		Message: "Session details retrieved",
	})
}

// HandleCreateSession_v0 verify user's credentials and generates a JWT token if valid
func HandleCreateSession_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	ta := ctx.GetTokenAuthorityCtx(req)

	var body CreateSessionRequest_v0
	err := controllers.DecodeJSONBodyHelper(w, req, r, &body)
	if err != nil {
		return
	}

	if body.Username == "" {
		controllers.RenderErrorResponseHelper(w, req, r, "Empty username")
		return
	}

	repo := repositories.UserRepository{db}
	user, err := repo.GetUserByUsername(body.Username)
	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	if !user.IsCredentialsVerified(body.Password) {
		controllers.RenderErrorResponseHelper(w, req, r, "Invalid username/password")
		return
	}

	var rolesString []string
	for _, role := range user.Roles {
		rolesString = append(rolesString, string(role))
	}

	tokenString, err := ta.CreateNewSessionToken(&domain.TokenClaims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Status:   user.Status,
		Roles:    rolesString,
	})

	if err != nil {
		controllers.RenderErrorResponseHelper(w, req, r, "Error creating session token")
		return
	}

	// run a post-create-session hook{
	hooks := ctx.GetControllerHooksMapCtx(req)
	if hooks.PostCreateSessionHook != nil {
		err = hooks.PostCreateSessionHook(w, req, ctx, &domain.PostCreateSessionHookPayload{
			TokenString: tokenString,
		})
		if err != nil {
			controllers.RenderErrorResponseHelper(w, req, r, err.Error())
			return
		}
	}
	// TODO: update user object with last logged-in

	r.JSON(w, http.StatusCreated, CreateSessionResponse_v0{
		Token:   tokenString,
		Success: true,
		Message: "Session token created",
	})
}

// HandleDeleteSession_v0 invalidates a session token
func HandleDeleteSession_v0(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	claims := ctx.GetAuthenticatedClaimsCtx(req)
	hooks := ctx.GetControllerHooksMapCtx(req)

	if claims == nil || !bson.IsObjectIdHex(claims.JTI) {
		// run a post-delete-session hook (
		if hooks.PostDeleteSessionHook != nil {
			err := hooks.PostDeleteSessionHook(w, req, ctx, &domain.PostDeleteSessionHookPayload{
				Claims: claims,
			})
			if err != nil {
				controllers.RenderErrorResponseHelper(w, req, r, err.Error())
				return
			}
		}
		// simply return because we can't blacklist a token without identifier
		r.JSON(w, http.StatusOK, DeleteSessionResponse_v0{
			Success: true,
			Message: "Session removed",
		})
		return
	}
	repo := repositories.RevokedTokenRepository{db}
	err := repo.CreateRevokedToken(&domain.RevokedToken{
		ID:         bson.ObjectIdHex(claims.JTI),
		ExpiryDate: claims.ExpireAt,
	})
	if err != nil {
		log.Println("HandleDeleteSession_v0: Failed to create revoked token", err.Error())
	}

	// run a post-delete-session hook{
	if hooks.PostDeleteSessionHook != nil {
		err = hooks.PostDeleteSessionHook(w, req, ctx, &domain.PostDeleteSessionHookPayload{
			Claims: claims,
		})
		if err != nil {
			controllers.RenderErrorResponseHelper(w, req, r, err.Error())
			return
		}
	}

	r.JSON(w, http.StatusOK, DeleteSessionResponse_v0{
		Success: true,
		Message: "Session removed",
	})
}
