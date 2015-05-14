package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"net/http"
	"strings"
)

type Authenticator struct {
}

func NewAuthenticator() *Authenticator {
	return &Authenticator{}
}

func (auth *Authenticator) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {

	r := ctx.GetRendererCtx(req)
	db := ctx.GetDbCtx(req)
	ta := ctx.GetTokenAuthorityCtx(req)

	authHeaderString := req.Header.Get("Authorization")
	if authHeaderString != "" {
		tokens := strings.Split(authHeaderString, " ")
		if len(tokens) != 2 || (len(tokens) > 0 && strings.ToUpper(tokens[0]) != "BEARER") {
			r.JSON(w, http.StatusUnauthorized, controllers.ErrorResponse_v0{
				Message: "Invalid format, expected Authorization: Bearer [token]",
				Success: false,
			})
			return
		}
		tokenString := tokens[1]
		token, claims, err := ta.VerifyTokenString(tokenString)
		if err != nil {
			r.JSON(w, http.StatusUnauthorized, controllers.ErrorResponse_v0{
				Message: "Unable to verify token string",
				Success: false,
			})
			return
		}
		if !token.Valid {
			r.JSON(w, http.StatusUnauthorized, controllers.ErrorResponse_v0{
				Message: "Token is invalid",
				Success: false,
			})
			return
		}

		// Check that the token was not previously revoked
		// TODO: Possible optimization, use Redis
		revokedTokenRepo := repositories.RevokedTokenRepository{db}
		if revokedTokenRepo.IsTokenRevoked(claims.JTI) {
			r.JSON(w, http.StatusUnauthorized, controllers.ErrorResponse_v0{
				Message: "Token has been revoked",
				Success: false,
			})
			return
		}

		// store claims for request context
		ctx.SetAuthenticatedClaimsCtx(req, claims)

		// retrieve user object and store it in current request context
		// this `user` object will be used by the AccessController middleware
		userRepo := repositories.UserRepository{db}
		user, _ := userRepo.GetUserById(claims.UserID)
		ctx.SetCurrentUserCtx(req, user)

	}
	next(w, req)
}
