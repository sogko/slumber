package domain

import (
	"net/http"
)

type IContext interface {
	SetDbCtx(r *http.Request, val IDatabase) IDatabase
	GetDbCtx(r *http.Request) IDatabase
	SetRendererCtx(r *http.Request, val IRenderer) IRenderer
	GetRendererCtx(r *http.Request) IRenderer
	SetTokenAuthorityCtx(r *http.Request, val ITokenAuthority) ITokenAuthority
	GetTokenAuthorityCtx(r *http.Request) ITokenAuthority
	SetAuthenticatedClaimsCtx(r *http.Request, val *TokenClaims) *TokenClaims
	GetAuthenticatedClaimsCtx(r *http.Request) *TokenClaims

	InjectWithNext(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc, ctx IContext)) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Inject(func(rw http.ResponseWriter, r *http.Request, ctx IContext)) http.HandlerFunc
}
