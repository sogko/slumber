package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/repositories"
	"log"
	"net/http"
	"strings"
)

type Authenticator struct {
}

func (auth *Authenticator) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {

	log.Println("Auth Handler")
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
		// store claims for request context
		ctx.SetAuthenticatedClaimsCtx(req, claims)

		repo := repositories.UserRepository{db}
		user, _ := repo.GetUserById(claims.UserID)
		log.Println("user", user)

		// retrieve user object
	}
	next(w, req)
}
