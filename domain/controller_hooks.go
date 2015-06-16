package domain

import (
	"net/http"
)

type ControllerHook func(w http.ResponseWriter, req *http.Request, ctx IContext, payload interface{}) error

//
//type ControllerHooksMap struct {
//	PostCreateUserHook    ControllerHook
//	PostConfirmUserHook   ControllerHook
//	PostCreateSessionHook ControllerHook
//	PostDeleteSessionHook ControllerHook
//}
//
//type PostCreateUserHookPayload struct {
//	User IUser
//}
//
//type PostConfirmUserHookPayload struct {
//	User IUser
//}
//
//type PostCreateSessionHookPayload struct {
//	TokenString string
//}
//
//type PostDeleteSessionHookPayload struct {
//	Claims interface{} // TODO: PostDeleteSessionHookPayload
////	Claims *TokenClaims
//}
