package server

import (
	"github.com/codegangsta/negroni"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

/*
 Database options

 Defaults:

 - ServerName   = ""

 - DatabaseName = ""

 - DialTimeout  = 60 seconds
*/
type DatabaseOptions struct {
	ServerName   string
	DatabaseName string
	DialTimeout  time.Duration
}

type DatabaseSession struct {
	*mgo.Session
	DatabaseOptions
}

// Returns a new database session
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

// Returns a negroni middleware HandlerFunc that creates and saves a database session into request context
func (session *DatabaseSession) UseDatabase() negroni.HandlerFunc {
	return negroni.HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// clone the `global` mgo session and save the named database in the request context for thread-safety
		s := session.Clone()
		defer s.Close()
		db := s.DB(session.DatabaseName)
		SetDbCtx(r, db)
		next(rw, r)
	})
}
