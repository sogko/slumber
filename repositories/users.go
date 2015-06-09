package repositories

import (
	"errors"
	"fmt"
	"github.com/sogko/slumber/domain"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

// User collection name
const UsersCollection string = "users"

type UserRepository struct {
	DB domain.IDatabase
}

// CreateUser Insert new user document into the database
func (repo *UserRepository) CreateUser(user *domain.User) error {
	user.ID = bson.NewObjectId()
	user.CreatedDate = time.Now()
	user.LastModifiedDate = time.Now()
	return repo.DB.Insert(UsersCollection, user)
}

// GetUsers Get list of users
func (repo *UserRepository) GetUsers() domain.Users {
	users := domain.Users{}
	err := repo.DB.FindAll(UsersCollection, nil, &users, 50, "")
	if err != nil {
		return domain.Users{}
	}
	return users
}

func (repo *UserRepository) FilterUsers(field string, query string, lastID string, limit int, sort string) domain.Users {
	users := domain.Users{}

	// ensure that collection has the right text index
	// refactor building collection index
	err := repo.DB.EnsureIndex(UsersCollection, mgo.Index{
		Key: []string{
			"$text:username",
			"$text:email",
			"$text:status",
		},
		Background: true,
		Sparse:     true,
	})
	if err != nil {
		log.Println("FilterUsers: EnsureIndex", err.Error())
	}
	// parse sort string
	allowedSortMap := map[string]bool{
		"_id":  true,
		"-_id": true,
	}
	// ensure that sort string is allowed
	// we are basically concerned about sorting on un-indexed keys
	if !allowedSortMap[sort] {
		sort = "-_id" // set it to default sort
	}

	q := domain.Query{}
	if lastID != "" && bson.IsObjectIdHex(lastID) {
		if sort == "_id" {
			q["_id"] = domain.Query{
				"$gt": bson.ObjectIdHex(lastID),
			}
		} else {
			q["_id"] = domain.Query{
				"$lt": bson.ObjectIdHex(lastID),
			}
		}
	}

	if query != "" {
		if field != "" {
			q[field] = domain.Query{
				"$regex":   fmt.Sprintf("^%v.*", query),
				"$options": "i",
			}
		} else {
			// if not field is specified, we do a text search on pre-defined text index
			q["$text"] = domain.Query{
				"$search": query,
			}
		}
	}

	err = repo.DB.FindAll(UsersCollection, q, &users, limit, sort)
	if err != nil {
		return domain.Users{}
	}
	return users
}

func (repo *UserRepository) CountUsers(field string, query string) int {
	q := domain.Query{}
	if query != "" {
		if field != "" {
			q[field] = domain.Query{
				"$regex":   fmt.Sprintf("^%v.*", query),
				"$options": "i",
			}
		} else {
			// if not field is specified, we do a text search on pre-defined text index
			q["$text"] = domain.Query{
				"$search": query,
			}
		}
	}

	count, err := repo.DB.Count(UsersCollection, q)
	if err != nil {
		return 0
	}
	return count
}

// DeleteUsers Delete a list of users
func (repo *UserRepository) DeleteUsers(ids []string) error {
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
	err := repo.DB.RemoveAll(UsersCollection, domain.Query{"_id": bson.M{"$in": objectIds}})
	return err
}

// DeleteAllUsers Delete all users
func (repo *UserRepository) DeleteAllUsers() error {
	err := repo.DB.DropCollection(UsersCollection)
	return err
}

// GetUser Get user specified by the id
func (repo *UserRepository) GetUserById(id string) (*domain.User, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	var user domain.User
	err := repo.DB.FindOne(UsersCollection, domain.Query{"_id": bson.ObjectIdHex(id)}, &user)
	return &user, err
}

// GetUser Get user specified by the username
func (repo *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := repo.DB.FindOne(UsersCollection, domain.Query{"username": username}, &user)
	return &user, err
}

// UserExistsByUsername Check if username already exists
func (repo *UserRepository) UserExistsByUsername(username string) bool {
	return repo.DB.Exists(UsersCollection, domain.Query{"username": username})
}

// UserExistsByEmail Check if email already exists
func (repo *UserRepository) UserExistsByEmail(email string) bool {
	return repo.DB.Exists(UsersCollection, domain.Query{"email": email})
}

// UpdateUser Update user specified by the id
func (repo *UserRepository) UpdateUser(id string, inUser *domain.User) (*domain.User, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}

	// serialize to a sub-set of allowed domain.User fields to update
	update := domain.Query{
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

	query := domain.Query{"_id": bson.ObjectIdHex(id)}
	change := domain.Change{
		Update:    domain.Query{"$set": update},
		ReturnNew: true,
	}
	var changedUser domain.User
	err := repo.DB.Update(UsersCollection, query, change, &changedUser)
	return &changedUser, err
}

// DeleteUser deletes user specified by the id
func (repo *UserRepository) DeleteUser(id string) error {

	if !bson.IsObjectIdHex(id) {
		return errors.New(fmt.Sprintf("Invalid ObjectId: `%v`", id))
	}
	err := repo.DB.RemoveOne(UsersCollection, domain.Query{"_id": bson.ObjectIdHex(id)})
	return err
}
