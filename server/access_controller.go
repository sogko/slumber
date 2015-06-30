package server

import (
	//	"github.com/sogko/slumber/controllers"
	"github.com/sogko/slumber/domain"
	"net/http"
)

const defaultForbiddenAccessMessage = "Forbidden (403)"
const defaultOKAccessMessage = "OK"

type ErrorResponse struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// TODO: Currently, AccessController only acts as a gateway for endpoints on router level. Build AC to handler other aspects of ACL
func NewAccessController(ctx domain.IContext, renderer domain.IRenderer) *AccessController {
	return &AccessController{domain.ACLMap{}, ctx, renderer}
}

// implements IAccessController
type AccessController struct {
	ACLMap   domain.ACLMap
	ctx      domain.IContext
	renderer domain.IRenderer
}

func (ac *AccessController) Add(aclMap *domain.ACLMap) {
	ac.ACLMap = ac.ACLMap.Append(aclMap)
}

func (ac *AccessController) AddHandler(action string, handler domain.ACLHandlerFunc) {
	ac.ACLMap[action] = handler
}

func (ac *AccessController) HasAction(action string) bool {
	fn := ac.ACLMap[action]
	return (fn != nil)
}

func (ac *AccessController) IsHTTPRequestAuthorized(req *http.Request, ctx domain.IContext, action string, user domain.IUser) (bool, string) {
	fn := ac.ACLMap[action]
	if fn == nil {
		// by default, if acl action/handler is not defined, request is not authorized
		return false, defaultForbiddenAccessMessage
	}

	result, message := fn(req, user)
	if result && message == "" {
		message = defaultOKAccessMessage
	}
	if !result && message == "" {
		message = defaultForbiddenAccessMessage
	}
	return result, message
}

func (ac *AccessController) NewContextHandler(action string, next http.HandlerFunc) http.HandlerFunc {
	//func (ac *AccessController) NewHandler(action string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := ac.ctx.GetCurrentUserCtx(req)
		// `user` might be `nil` if has not authenticated.
		// ACL might want to allow anonymous / non-authenticated access (for login, e.g)

		result, message := ac.IsHTTPRequestAuthorized(req, ac.ctx, action, user)
		if !result {
			ac.renderer.Render(w, req, http.StatusForbidden, ErrorResponse{
				Message: message,
				Success: false,
			})
			return
		}

		next(w, req)
	}
}
