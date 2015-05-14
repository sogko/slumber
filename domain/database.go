package domain

import (
	"gopkg.in/mgo.v2"
	"net/http"
)

type IDatabaseOptions interface {
}

type Query map[string]interface{}
type Change mgo.Change

// Database interface
type IDatabase interface {
	NewSession() IDatabaseSession

	Insert(name string, obj interface{}) error
	Update(name string, query Query, change Change, result interface{}) error
	FindOne(name string, query Query, result interface{}) error
	FindAll(name string, query Query, result interface{}) error
	Remove(name string, query Query) error
	RemoveAll(name string) error
	Exists(name string, query Query) bool
	DropDatabase() error
}

type IDatabaseSession interface {
	Handler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)
}
