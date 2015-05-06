package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

var privateSigningKey []byte
var publicSigningKey []byte

func init() {
	privateSigningKey, _ = ioutil.ReadFile("../keys/demo.rsa")
	publicSigningKey, _ = ioutil.ReadFile("../keys/demo.rsa.pub")
}

type TokenClaims struct {
	ID     string
	Status string
	Roles  []string
}

func CreateNewToken(claims TokenClaims) (string, error) {

	token := jwt.New(jwt.SigningMethodRS512)

	token.Claims = map[string]interface{}{
		"id":     claims.ID,
		"status": claims.Status,
		"roles":  claims.Roles,
		"exp":    time.Now().Add(time.Hour * 72).Unix(), // 3 days
		"iat":    time.Now(),
		"jti":    "nonce",
	}

	tokenString, err := token.SignedString(privateSigningKey)

	return tokenString, err
}

func ParseAndVerifyTokenString(tokenStr string) (*jwt.Token, *TokenClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicSigningKey, nil
	})
	if err != nil {
		return nil, nil, err
	}

	var claims TokenClaims
	if token.Valid {
		claims.ID = token.Claims["id"].(string)
		claims.Status = token.Claims["status"].(string)
		for _, role := range token.Claims["roles"].([]interface{}) {
			claims.Roles = append(claims.Roles, role.(string))
		}
	}

	return token, &claims, err
}
