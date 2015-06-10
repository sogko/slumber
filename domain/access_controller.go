package domain

import (
	"net/http"
)

type ACLHandlerFunc func(*User, *http.Request, IContext) (bool, string)
type ACLMap map[string]ACLHandlerFunc

type IAccessController interface {
	Add(*ACLMap)
	HasAction(string) bool
	IsHTTPRequestAuthorized(req *http.Request, ctx IContext, action string, user *User) (bool, string)
	Handler(action string, handler ContextHandlerFunc) ContextHandlerFunc
}
