package domain

import (
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

type ITokenClaims interface {
}
type TokenClaims struct {
	UserID   string
	Username string
	Status   string
	Roles    []string
	ExpireAt time.Time
	IssuedAt time.Time
	JTI      string
}
type RevokedToken struct {
	ID          bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	ExpiryDate  time.Time     `json:"exp" bson:"exp"`
	RevokedDate time.Time     `json:"revokedat" bson:"revokedat"`
}
type Token jwt.Token
type ITokenAuthorityOptions interface{}

type ITokenAuthority interface {
	CreateNewSessionToken(claims *TokenClaims) (string, error)
	VerifyTokenString(tokenString string) (Token, *TokenClaims, error)
	Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx IContext)
}
