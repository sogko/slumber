package middlewares

import (
	"github.com/sogko/slumber/domain"
	"net/http"
)

type ControllerHooksMiddleware struct {
	ControllerHooksMap *domain.ControllerHooksMap
}

func (c *ControllerHooksMiddleware) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	ctx.SetControllerHooksMapCtx(req, c.ControllerHooksMap)
	next(w, req)
}
