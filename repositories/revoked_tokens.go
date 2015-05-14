package repositories

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User collection name
const RevokedTokenCollections string = "revoked_tokens"

type RevokedTokenRepository struct {
	DB domain.IDatabase
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) CreateRevokedToken(token *domain.RevokedToken) error {
	token.RevokedDate = time.Now()
	return repo.DB.Insert(RevokedTokenCollections, token)
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) DeleteExpiredTokens() error {
	return repo.DB.RemoveAll(RevokedTokenCollections, domain.Query{
		"exp": domain.Query{
			"$lt": time.Now(),
		},
	})
}

// CreateRevokedToken Insert new user document into the database
func (repo *RevokedTokenRepository) IsTokenRevoked(id string) bool {
	if !bson.IsObjectIdHex(id) {
		return false
	}
	return repo.DB.Exists(RevokedTokenCollections, domain.Query{
		"_id": bson.ObjectIdHex(id),
	})
}
