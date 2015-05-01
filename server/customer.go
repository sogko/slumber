package server

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const CustomersCollection = "customers"

type Customer struct {
	Id               bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	BusinessFileId   string        `json:"businessFileId" bson:"businessFileId,omitempty"`
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

type Customers []Customer

// Ensure that the customer object is valid
func (customer *Customer) IsValid() bool {
	return len(customer.FirstName) > 0 &&
		len(customer.Email) > 0
}

// Insert new customer document into the database
func CreateCustomer(db *mgo.Database, customer *Customer) error {
	customer.Id = bson.NewObjectId()
	return db.C(CustomersCollection).Insert(customer)
}

// Get list of customers
func GetCustomers(db *mgo.Database) (Customers, error) {
	customers := Customers{}
	err := db.C(CustomersCollection).Find(nil).All(&customers)
	return customers, err
}

// Get customer specified by the id
func GetCustomer(db *mgo.Database, id string) (*Customer, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid ObjectId")
	}

	var customer Customer
	err := db.C(CustomersCollection).FindId(bson.ObjectIdHex(id)).One(&customer)
	return &customer, err
}
