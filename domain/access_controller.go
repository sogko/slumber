package domain

import (
	"net/http"
)

type ACLHandlerFunc func(*User, *http.Request, IContext) (bool, string)

type ACLMap map[string]ACLHandlerFunc

func (m *ACLMap) Append(maps ...*ACLMap) ACLMap {
	res := ACLMap{}
	// copy current map
	for k, v := range *m {
		res[k] = v
	}
	for _, _maps := range maps {
		for k, v := range *_maps {
			res[k] = v
		}
	}
	return res
}

type IAccessController interface {
	Add(*ACLMap)
	HasAction(string) bool
	IsHTTPRequestAuthorized(req *http.Request, ctx IContext, action string, user *User) (bool, string)
	Handler(action string, handler ContextHandlerFunc) ContextHandlerFunc
}
