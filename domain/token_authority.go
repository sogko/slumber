package domain

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

type ITokenClaims interface {
}
type TokenClaims struct {
	UserID   string
	Status   string
	Roles    []string
	ExpireAt time.Time
	IssuedAt time.Time
	JTI      string
}
type Token jwt.Token
type ITokenAuthorityOptions interface{}

type ITokenAuthority interface {
	CreateNewSessionToken(claims *TokenClaims) (string, error)
	VerifyTokenString(tokenString string) (Token, *TokenClaims, error)
	Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx IContext)
}
