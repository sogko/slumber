package middlewares

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func generateJTI() string {
	// We will use mongodb's object id as JTI
	// we then will use this id to blacklist tokens,
	// along with `exp` and `iat` claims.
	// As far as collisions go, ObjectId is guaranteed unique
	// within a collection; and this case our collection is `sessions`
	return bson.NewObjectId().Hex()
}

func NewTokenAuthority(options domain.ITokenAuthorityOptions) *TokenAuthority {
	ta := TokenAuthority{}
	if options != nil {
		ta.Options = *options.(*TokenAuthorityOptions)
	}
	return &ta
}

type TokenAuthority struct {
	Options TokenAuthorityOptions
}

type TokenAuthorityOptions struct {
	PrivateSigningKey []byte
	PublicSigningKey  []byte
}

func (ta *TokenAuthority) CreateNewSessionToken(claims *domain.TokenClaims) (string, error) {

	token := jwt.New(jwt.SigningMethodRS512)

	token.Claims = map[string]interface{}{
		"user_id": claims.UserID,
		"status":  claims.Status,
		"roles":   claims.Roles,
		"exp":     time.Now().Add(time.Hour * 72).Format(time.RFC3339), // 3 days
		"iat":     time.Now().Format(time.RFC3339),
		"jti":     generateJTI(),
	}
	tokenString, err := token.SignedString(ta.Options.PrivateSigningKey)

	return tokenString, err
}

func (ta *TokenAuthority) VerifyTokenString(tokenStr string) (domain.Token, *domain.TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return ta.Options.PublicSigningKey, nil
	})
	if err != nil {
		return domain.Token{}, nil, err
	}

	var claims domain.TokenClaims
	if token.Valid {
		if token.Claims["user_id"] != nil {
			claims.UserID = token.Claims["user_id"].(string)
		}
		if token.Claims["status"] != nil {
			claims.Status = token.Claims["status"].(string)
		}
		if token.Claims["roles"] != nil {
			for _, role := range token.Claims["roles"].([]interface{}) {
				if role != nil {
					claims.Roles = append(claims.Roles, role.(string))
				}
			}
		}
		if token.Claims["jti"] != nil {
			claims.JTI = token.Claims["jti"].(string)
		}
		if token.Claims["iat"] != nil {
			claims.IssuedAt, _ = time.Parse(time.RFC3339, token.Claims["iat"].(string))
		}
		if token.Claims["exp"] != nil {
			claims.ExpireAt, _ = time.Parse(time.RFC3339, token.Claims["exp"].(string))
		}
	}

	return domain.Token(*token), &claims, err
}
func (ta *TokenAuthority) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	ctx.SetTokenAuthorityCtx(req, ta)
	next(w, req)

}
