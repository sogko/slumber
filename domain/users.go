package domain

import (
	"github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

// User statuses
const (
	StatusPending   = "pending"
	StatusActive    = "active"
	StatusInactive  = "inactive"
	StatusSuspended = "suspended"
	StatusDeleted   = "deleted"
)

type NewUser struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

// User model struct
// User refers to the REST API user.
type User struct {
	ID               bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Username         string        `json:"username,omitempty" bson:"username"`
	Email            string        `json:"email,omitempty" bson:"email"`
	Roles            Roles         `json:"roles,omitempty" bson:"roles"`
	Status           string        `json:"status,omitempty" bson:"status"`
	LastModifiedDate time.Time     `json:"lastModifiedDate" bson:"lastModifiedDate"`
	CreatedDate      time.Time     `json:"createdDate,omitempty" bson:"createdDate"`

	// fields are not exported to JSON
	ConfirmationCode string `json:"-" bson:"confirmationCode"`
	HashedPassword   string `json:"-" bson:"hashedPassword"`
}

// Users struct
type Users []User

// IsValid Ensures that the customer object is valid
func (user *User) IsValid() bool {
	// Note regarding email address validation,
	// as long as it `*looks* like an address, we'll allow it.
	// User need to confirm account by click on link sent to the email anyways
	return len(user.Username) > 0 &&
		len(user.Email) > 0 &&
		strings.Contains(user.Email, "@")
}

// IsCodeVerified verify the given code
func (user *User) IsCodeVerified(code string) bool {
	return (user.ConfirmationCode == code)
}

// IsCredentialsVerified verify the given credentials
func (user *User) IsCredentialsVerified(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return (err == nil)
}

// SetPassword encrypts the given plain text password
func (user *User) SetPassword(password string) error {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.HashedPassword = string(hashedPassword)
	return nil
}

// SetPassword encrypts the given plain text password
func (user *User) GenerateConfirmationCode() {
	user.ConfirmationCode = generateNewUniqueCode()
}

// generateNewUniqueCode generates a new confirmation code
func generateNewUniqueCode() string {
	// set code format
	uuid.SwitchFormat(uuid.Clean)
	return uuid.NewV4().String()
}

func (user *User) HasRole(r Role) bool {
	for _, a := range user.Roles {
		if a == r {
			return true
		}
	}
	return false
}
