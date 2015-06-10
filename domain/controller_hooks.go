package domain

import (
	"net/http"
)

type ControllerHook func(w http.ResponseWriter, req *http.Request, ctx IContext, payload interface{}) error
type ControllerHooksMap struct {
	PostCreateUserHook    ControllerHook
	PostConfirmUserHook   ControllerHook
	PostCreateSessionHook ControllerHook
	PostDeleteSessionHook ControllerHook
}

type PostCreateUserHookPayload struct {
	User *User
}

type PostConfirmUserHookPayload struct {
	User *User
}

type PostCreateSessionHookPayload struct {
	TokenString string
}

type PostDeleteSessionHookPayload struct {
	Claims *TokenClaims
}
