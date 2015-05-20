package domain

import (
	"net/http"
)

type ControllerHook func(req *http.Request, ctx IContext, payload interface{}) error
type ControllerHooksMap struct {
	PostConfirmUserHook ControllerHook
}

type PostUserConfirmationHookPayload struct {
	User *User
}
