package middlewares

import (
	"github.com/gorilla/context"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"net/http"
)

type contextKey string

const (
	RouteKey          contextKey = "RouteKey"
	DbKey             contextKey = "DbKey"
	RendererKey       contextKey = "RendererKey"
	TokenClaimsKey    contextKey = "TokenClaimsKey"
	TokenAuthorityKey contextKey = "TokenAuthorityKey"
	CurrentUserKey    contextKey = "CurrentUserKey"
	CurrentObjectKey  contextKey = "CurrentObjectKey"
)

func NewContext() *Context {
	return &Context{}
}

// implements domain.IContext
type Context struct {
}

func (ctx *Context) InjectWithNext(middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx domain.IContext)) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

// SetRouteCtx Sets the Database reference for the given request context
func (ctx *Context) SetRouteCtx(r *http.Request, val *domain.Route) *domain.Route {
	context.Set(r, RouteKey, val)
	return val
}

// GetRouteCtx Returns the Database reference for the given request context
func (ctx *Context) GetRouteCtx(r *http.Request) *domain.Route {
	if r := context.Get(r, RouteKey); r != nil {
		return r.(*domain.Route)
	}
	return &domain.Route{}
}

// SetDbCtx Sets the Database reference for the given request context
func (ctx *Context) SetDbCtx(r *http.Request, val domain.IDatabase) domain.IDatabase {
	context.Set(r, DbKey, val)
	return val
}

// DbCtx Returns the Database reference for the given request context
func (ctx *Context) GetDbCtx(r *http.Request) domain.IDatabase {
	if db := context.Get(r, DbKey); db != nil {
		return db.(domain.IDatabase)
	}
	return nil
}

// SetRendererCtx Set the Render reference for the given request context
func (ctx *Context) SetRendererCtx(r *http.Request, val domain.IRenderer) domain.IRenderer {
	context.Set(r, RendererKey, val)
	return val
}

// RendererCtx Returns the Render reference for the given request context
func (ctx *Context) GetRendererCtx(r *http.Request) domain.IRenderer {
	if r := context.Get(r, RendererKey); r != nil {
		return r.(domain.IRenderer)
	}
	return nil
}

// SetAuthenticatedClaimsCtx Set the TokenClaims reference for the given request context
func (ctx *Context) SetAuthenticatedClaimsCtx(r *http.Request, val *domain.TokenClaims) *domain.TokenClaims {
	context.Set(r, TokenClaimsKey, val)
	return val
}

// GetAuthenticatedClaimsCtx the TokenClaims reference for the given request context
func (ctx *Context) GetAuthenticatedClaimsCtx(r *http.Request) *domain.TokenClaims {
	if r := context.Get(r, TokenClaimsKey); r != nil {
		return r.(*domain.TokenClaims)
	}
	return nil
}

// SetTokenAuthorityCtx Set the TokenAuthority reference for the given request context
func (ctx *Context) SetTokenAuthorityCtx(r *http.Request, val domain.ITokenAuthority) domain.ITokenAuthority {
	context.Set(r, TokenAuthorityKey, val)
	return val
}

// GetTokenAuthorityCtx the TokenAuthority reference for the given request context
func (ctx *Context) GetTokenAuthorityCtx(r *http.Request) domain.ITokenAuthority {
	if r := context.Get(r, TokenAuthorityKey); r != nil {
		return r.(domain.ITokenAuthority)
	}
	return nil
}

// SetCurrentUserCtx Set the TokenAuthority reference for the given request context
func (ctx *Context) SetCurrentUserCtx(r *http.Request, val *domain.User) *domain.User {
	context.Set(r, CurrentUserKey, val)
	return val
}

// GetCurrentUserCtx the TokenAuthority reference for the given request context
func (ctx *Context) GetCurrentUserCtx(r *http.Request) *domain.User {
	if r := context.Get(r, CurrentUserKey); r != nil {
		return r.(*domain.User)
	}
	return nil
}

// SetCurrentObjectCtx Set the TokenAuthority reference for the given request context
func (ctx *Context) SetCurrentObjectCtx(r *http.Request, val interface{}) interface{} {
	context.Set(r, CurrentObjectKey, val)
	return val
}

// GetCurrentObjectCtx the TokenAuthority reference for the given request context
func (ctx *Context) GetCurrentObjectCtx(r *http.Request) interface{} {
	if r := context.Get(r, CurrentObjectKey); r != nil {
		return r
	}
	return nil
}
