package models

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Customer model struct
type Customer struct {
	ID               bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	BusinessFileID   string        `json:"businessFileId" bson:"businessFileId,omitempty"`
	Owner            bson.ObjectId `json:"owner" bson:"owner,omitempty"`
	FirstName        string        `json:"firstName"`
	LastName         string        `json:"lastName"`
	Email            string        `json:"email"`
	Notes            string        `json:"notes"`
	Currency         string        `json:"currency"`
	Location         string        `json:"location"`
	Phone            string        `json:"phone"`
	LastModifiedDate time.Time     `json:"lastModifiedDate"`
	CreatedDate      time.Time     `json:"createdDate"`
}

// Customers struct
type Customers []Customer

// IsValid Ensures that the customer object is valid
func (customer *Customer) IsValid() bool {
	return len(customer.FirstName) > 0 &&
		len(customer.Email) > 0
}
