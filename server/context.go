package server

import (
	"github.com/gorilla/context"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"net/http"
)

type key int

// Request context keys
const DbKey key = 0
const RendererKey key = 1

// Sets the mgo.Database reference for the given request context
func SetDbCtx(r *http.Request, val *mgo.Database) {
	context.Set(r, DbKey, val)
}

// Returns the mgo.Database reference for the given request context
func DbCtx(r *http.Request) *mgo.Database {
	if db := context.Get(r, DbKey); db != nil {
		return db.(*mgo.Database)
	}
	return nil
}

// Set the render.Render reference for the given request context
func SetRenderCtx(r *http.Request, val *render.Render) {
	context.Set(r, RendererKey, val)
}

// Returns the render.Render reference for the given request context
func RenderCtx(r *http.Request) *render.Render {
	if r := context.Get(r, RendererKey); r != nil {
		return r.(*render.Render)
	}
	return nil
}
