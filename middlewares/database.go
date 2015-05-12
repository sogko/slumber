package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/domain"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type MongoDBOptions struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

// MongoDatabase implements Database interface
type MongoDB struct {
	currentDb *mgo.Database
}

// CreateSession Returns a new database session
// Defaults DatabaseOptions:
// - ServerName   = ""
// - DatabaseName = ""
// - DialTimeout  = 60 seconds
//
func (db *MongoDB) NewSession(options domain.IDatabaseOptions) domain.IDatabaseSession {

	var mongoOptions = options.(*MongoDBOptions)

	// set default DialTimeout value
	if mongoOptions.DialTimeout <= 0 {
		mongoOptions.DialTimeout = 1 * time.Minute
	}

	session, err := mgo.DialWithTimeout(mongoOptions.ServerName, mongoOptions.DialTimeout)
	if err != nil {
		panic(err)
	}
	return &MongoDBSession{session, mongoOptions}
}

func (db *MongoDB) FindOne(name string, query domain.Query, result interface{}) error {
	return db.currentDb.C(name).Find(query).One(result)
}

func (db *MongoDB) FindAll(name string, query domain.Query, result interface{}) error {
	return db.currentDb.C(name).Find(query).All(result)
}

func (db *MongoDB) Insert(name string, obj interface{}) error {
	return db.currentDb.C(name).Insert(obj)
}

func (db *MongoDB) Update(name string, query domain.Query, change domain.Change, result interface{}) error {
	_, err := db.currentDb.C(name).Find(query).Apply(mgo.Change(change), result)
	return err
}

func (db *MongoDB) Remove(name string, query domain.Query) error {
	return db.currentDb.C(name).Remove(query)
}

func (db *MongoDB) RemoveAll(name string) error {
	return db.currentDb.C(name).DropCollection()
}

func (db *MongoDB) Exists(name string, query domain.Query) bool {
	var result interface{}
	err := db.currentDb.C(name).Find(query).One(result)
	return (err == nil)
}

// MongoDatabaseSession struct implements DatabaseSession interface
type MongoDBSession struct {
	*mgo.Session
	*MongoDBOptions
}

// HandlerWithNext Returns a middleware HandlerFunc that creates and saves a database session into request context
func (session *MongoDBSession) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	// clone the `global` mgo session and save the named database in the request context for thread-safety
	s := session.Clone()
	defer s.Close()
	db := &MongoDB{s.DB(session.DatabaseName)}
	ctx.SetDbCtx(req, db)
	next(w, req)
}
