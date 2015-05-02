package server

import (
	"errors"
	"github.com/sogko/golang-rest-api-server-example/server/models"
	"gopkg.in/mgo.v2/bson"
)

// Customer collection name
const CustomersCollection string = "customers"

// CreateCustomer Insert new customer document into the database
func (db *Database) CreateCustomer(customer *models.Customer) error {
	customer.ID = bson.NewObjectId()
	return db.C(CustomersCollection).Insert(customer)
}

// GetCustomers Get list of customers
func (db *Database) GetCustomers() (models.Customers, error) {
	customers := models.Customers{}
	err := db.C(CustomersCollection).Find(nil).All(&customers)
	return customers, err
}

// GetCustomer Get customer specified by the id
func (db *Database) GetCustomer(id string) (*models.Customer, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid ObjectId")
	}

	var customer models.Customer
	err := db.C(CustomersCollection).FindId(bson.ObjectIdHex(id)).One(&customer)
	return &customer, err
}
