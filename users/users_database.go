package users

import (
	"errors"
	"fmt"
	"github.com/sogko/golang-rest-api-server-example/server"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// User collection name
const UsersCollection string = "users"

// CreateUser Insert new user document into the database
func CreateUser(db *server.Database, user *User) error {
	user.ID = bson.NewObjectId()
	user.CreatedDate = time.Now()
	user.LastModifiedDate = time.Now()
	return db.C(UsersCollection).Insert(user)
}

// GetUsers Get list of users
func GetUsers(db *server.Database) Users {
	users := Users{}
	err := db.C(UsersCollection).Find(nil).All(&users)
	if err != nil {
		return Users{}
	}
	return users
}

// DeleteUsers Delete a list of users
func DeleteUsers(db *server.Database, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	var objectIds []bson.ObjectId
	for _, id := range ids {
		if bson.IsObjectIdHex(id) {
			objectIds = append(objectIds, bson.ObjectIdHex(id))
		}
	}
	if len(objectIds) == 0 {
		return nil
	}
	err := db.C(UsersCollection).Remove(bson.M{"_id": bson.M{"$in": objectIds}})
	return err
}

// DeleteAllUsers Delete all users
func DeleteAllUsers(db *server.Database) error {
	err := db.C(UsersCollection).DropCollection()
	return err
}

// GetUser Get user specified by the id
func GetUser(db *server.Database, id string) (*User, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var user User
	err := db.C(UsersCollection).FindId(bson.ObjectIdHex(id)).One(&user)
	return &user, err
}

// UpdateUser Update user specified by the id
func UpdateUser(db *server.Database, id string, inUser *User) (*User, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	// serialize to a sub-set of allowed User fields to update
	update := bson.M{
		"lastModifiedDate": time.Now(),
	}
	if inUser.Email != "" {
		update["email"] = inUser.Email
	}
	if inUser.Username != "" {
		update["username"] = inUser.Username
	}
	if inUser.Status != "" {
		update["status"] = inUser.Status
	}
	if len(inUser.Roles) > 0 {
		update["roles"] = inUser.Roles
	}

	change := mgo.Change{
		Update:    bson.M{"$set": update},
		ReturnNew: true,
	}

	var changedUser User
	_, err := db.C(UsersCollection).Find(bson.M{"_id": bson.ObjectIdHex(id)}).Apply(change, &changedUser)

	return &changedUser, err
}

// DeleteUser deletes user specified by the id
func DeleteUser(db *server.Database, id string) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := db.C(UsersCollection).Remove(bson.M{"_id": bson.ObjectIdHex(id)})
	return err
}
