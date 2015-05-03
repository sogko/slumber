package customers

import (
	"errors"
	"github.com/sogko/golang-rest-api-server-example/server"
	"gopkg.in/mgo.v2/bson"
)

// Customer collection name
const CustomersCollection string = "customers"

// CreateCustomer Insert new customer document into the database
func CreateCustomer(db *server.Database, customer *Customer) error {
	customer.ID = bson.NewObjectId()
	return db.C(CustomersCollection).Insert(customer)
}

// GetCustomers Get list of customers
func GetCustomers(db *server.Database) (Customers, error) {
	customers := Customers{}
	err := db.C(CustomersCollection).Find(nil).All(&customers)
	return customers, err
}

// GetCustomer Get customer specified by the id
func GetCustomer(db *server.Database, id string) (*Customer, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid ObjectId")
	}

	var customer Customer
	err := db.C(CustomersCollection).FindId(bson.ObjectIdHex(id)).One(&customer)
	return &customer, err
}
