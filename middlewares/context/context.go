package context

import (
	"github.com/gorilla/context"
	"github.com/sogko/slumber/domain"
	"net/http"
)

const CurrentUserKey domain.ContextKey = "slumber-mddlwr-context-current-user-key"
const DatabaseKey domain.ContextKey = "slumber-mddlwr-context-database-key"

func New() *Context {
	return &Context{}
}

// Context implements IContext
type Context struct {
}

func (ctx *Context) InjectMiddleware(middleware domain.ContextMiddlewareFunc) domain.MiddlewareFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		middleware(rw, r, next, ctx)
	}
}

func (ctx *Context) Inject(handler domain.ContextHandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r, ctx)
	}
}

func (ctx *Context) Set(r *http.Request, key interface{}, val interface{}) {
	context.Set(r, key, val)
}

func (ctx *Context) Get(r *http.Request, key interface{}) interface{} {
	return context.Get(r, key)
}

func (ctx *Context) SetCurrentUserCtx(r *http.Request, user domain.IUser) {
	ctx.Set(r, CurrentUserKey, user)
}

func (ctx *Context) GetCurrentUserCtx(r *http.Request) domain.IUser {
	if user := ctx.Get(r, CurrentUserKey); user != nil {
		return user.(domain.IUser)
	}
	return nil
}
