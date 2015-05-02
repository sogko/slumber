package server

import (
	"errors"
	"github.com/codegangsta/negroni"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// Customer collection name
const CustomersCollection string = "customers"

// DatabaseOptions struct
//
// Defaults:
//
// - ServerName   = ""
//
// - DatabaseName = ""
//
// - DialTimeout  = 60 seconds
//
type DatabaseOptions struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

// DatabaseSession struct
type DatabaseSession struct {
	*mgo.Session
	DatabaseOptions
}

// Database type
type Database struct {
	*mgo.Database
}

// NewSession Returns a new database session
func NewSession(options DatabaseOptions) *DatabaseSession {

	// set default DialTimeout value
	if options.DialTimeout <= 0 {
		options.DialTimeout = 1 * time.Minute
	}

	session, err := mgo.DialWithTimeout(options.ServerName, options.DialTimeout)
	if err != nil {
		panic(err)
	}
	return &DatabaseSession{session, options}
}

// UseDatabase Returns a negroni middleware HandlerFunc that creates and saves a database session into request context
func (session *DatabaseSession) UseDatabase() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// clone the `global` mgo session and save the named database in the request context for thread-safety
		s := session.Clone()
		defer s.Close()
		db := &Database{s.DB(session.DatabaseName)}
		SetDbCtx(r, db)
		next(rw, r)
	})
}

//---- Customer database operations -----

// CreateCustomer Insert new customer document into the database
func (db *Database) CreateCustomer(customer *Customer) error {
	customer.ID = bson.NewObjectId()
	return db.C(CustomersCollection).Insert(customer)
}

// GetCustomers Get list of customers
func (db *Database) GetCustomers() (Customers, error) {
	customers := Customers{}
	err := db.C(CustomersCollection).Find(nil).All(&customers)
	return customers, err
}

// GetCustomer Get customer specified by the id
func (db *Database) GetCustomer(id string) (*Customer, error) {

	if !bson.IsObjectIdHex(id) {
		return nil, errors.New("Invalid ObjectId")
	}

	var customer Customer
	err := db.C(CustomersCollection).FindId(bson.ObjectIdHex(id)).One(&customer)
	return &customer, err
}
