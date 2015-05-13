package middlewares

import (
	"github.com/sogko/golang-rest-api-server-example/controllers"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"net/http"
)

func NewAccessController() *AccessController {
	ac := AccessController{}
	ac.ACLMap = domain.ACLMap{}
	return &ac
}

// MergeACLMap Returns a new map
func MergeACLMap(to *domain.ACLMap, from *domain.ACLMap) domain.ACLMap {
	res := domain.ACLMap{}
	for k, v := range *to {
		res[k] = v
	}
	for k, v := range *from {
		res[k] = v
	}
	return res
}

// implements IAccessController
type AccessController struct {
	req    *http.Request
	ctx    domain.IContext
	ACLMap domain.ACLMap
}

func (ac *AccessController) SetRequestContext(req *http.Request, ctx domain.IContext) {
	ac.req = req
	ac.ctx = ctx
}

func (ac *AccessController) Add(_aclMap *domain.ACLMap) {
	ac.ACLMap = MergeACLMap(&ac.ACLMap, _aclMap)
}

func (ac *AccessController) HasAction(action string) bool {
	fn := ac.ACLMap[action]
	return (fn != nil)
}

func (ac *AccessController) IsAuthorized(action string, user *domain.User) bool {
	fn := ac.ACLMap[action]
	if fn == nil {
		// by default, if acl action/handler is not defined, request is not authorized
		return false
	}
	return fn(user, ac.req, ac.ctx)
}

func (ac *AccessController) Handler(action string, handler domain.ContextHandlerFunc) domain.ContextHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
		r := ctx.GetRendererCtx(req)
		user := ctx.GetCurrentUserCtx(req)

		// `user` might be `nil` if has not authenticated.
		// ACL might want to allow anonymous / non-authenticated access (for login, e.g)

		if !ac.IsAuthorized(action, user) {
			r.JSON(w, http.StatusForbidden, controllers.ErrorResponse_v0{
				Message: "Forbidden (403)",
				Success: false,
			})
			return
		}

		handler(w, req, ctx)
	}
}
