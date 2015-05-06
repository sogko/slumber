package sessions

import (
	"github.com/sogko/golang-rest-api-server-example/server"
	"github.com/sogko/golang-rest-api-server-example/users"
	"gopkg.in/mgo.v2/bson"
)

// User collection name
const UsersCollection string = "users"

// GetUser Get user specified by the username
func GetUserByUsername(db *server.Database, username string) (*users.User, error) {
	var user users.User
	err := db.C(UsersCollection).Find(bson.M{"username": username}).One(&user)
	return &user, err
}
