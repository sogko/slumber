package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/libs"
	"net/http"
)

// TODO: Currently, AccessController only acts as a gateway for endpoints on router level. Build AC to handler other aspects of ACL
func NewAccessController() *AccessController {
	ac := AccessController{}
	ac.ACLMap = domain.ACLMap{}
	return &ac
}

// implements IAccessController
type AccessController struct {
	ACLMap domain.ACLMap
}

func (ac *AccessController) Add(_aclMap *domain.ACLMap) {
	ac.ACLMap = libs.MergeACLMap(&ac.ACLMap, _aclMap)
}

func (ac *AccessController) HasAction(action string) bool {
	fn := ac.ACLMap[action]
	return (fn != nil)
}

func (ac *AccessController) IsHTTPRequestAuthorized(req *http.Request, ctx domain.IContext, action string, user *domain.User) bool {
	fn := ac.ACLMap[action]
	if fn == nil {
		// by default, if acl action/handler is not defined, request is not authorized
		return false
	}
	return fn(user, req, ctx)
}

func (ac *AccessController) Handler(action string, handler domain.ContextHandlerFunc) domain.ContextHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
		r := ctx.GetRendererCtx(req)
		user := ctx.GetCurrentUserCtx(req)

		// `user` might be `nil` if has not authenticated.
		// ACL might want to allow anonymous / non-authenticated access (for login, e.g)

		if !ac.IsHTTPRequestAuthorized(req, ctx, action, user) {
			r.JSON(w, http.StatusForbidden, controllers.ErrorResponse_v0{
				Message: "Forbidden (403)",
				Success: false,
			})
			return
		}

		handler(w, req, ctx)
	}
}
