package server

import (
	"github.com/gorilla/context"
	"net/http"
)

type key int

// Request context key for Database
const DbKey key = 0

// Request context key for Renderer
const RendererKey key = 1

// SetDbCtx Sets the Database reference for the given request context
func SetDbCtx(r *http.Request, val *Database) *Database {
	context.Set(r, DbKey, val)
	return val
}

// DbCtx Returns the Database reference for the given request context
func DbCtx(r *http.Request) *Database {
	if db := context.Get(r, DbKey); db != nil {
		return db.(*Database)
	}
	return nil
}

// SetRendererCtx Set the Render reference for the given request context
func SetRendererCtx(r *http.Request, val *Renderer) *Renderer {
	context.Set(r, RendererKey, val)
	return val
}

// RendererCtx Returns the Render reference for the given request context
func RendererCtx(r *http.Request) *Renderer {
	if r := context.Get(r, RendererKey); r != nil {
		return r.(*Renderer)
	}
	return nil
}
