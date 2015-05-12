package middlewares

import (
	"github.com/gorilla/context"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"net/http"
)

type contextKey string

const (
	DbKey             contextKey = "DbKey"
	RendererKey       contextKey = "RendererKey"
	TokenClaimsKey    contextKey = "TokenClaimsKey"
	TokenAuthorityKey contextKey = "TokenAuthorityKey"
)

// implements domain.IContext
type Context struct {
}

func (ctx *Context) InjectWithNext(middleware func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx domain.IContext)) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		middleware(rw, r, next, ctx)
	}
}

func (ctx *Context) Inject(handler func(rw http.ResponseWriter, r *http.Request, ctx domain.IContext)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		handler(rw, r, ctx)
	}
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

// SetAuthenticatedClaimsCtx the TokenClaims reference for the given request context
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
