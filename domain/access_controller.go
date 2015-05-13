package domain

import (
	"net/http"
)

type ACLHandlerFunc func(*User, *http.Request, IContext) bool
type ACLMap map[string]ACLHandlerFunc

type IAccessController interface {
	Add(*ACLMap)
	HasAction(string) bool
	IsAuthorized(string, *User) bool
	Handler(action string, handler ContextHandlerFunc) ContextHandlerFunc
}
