package domain

import (
	"net/http"
)

// Renderer interface
type IRenderer interface {
	Render(w http.ResponseWriter, req *http.Request, status int, v interface{})
	JSON(w http.ResponseWriter, status int, v interface{})
	XML(w http.ResponseWriter, status int, v interface{})
	Data(w http.ResponseWriter, status int, v []byte)
	Text(w http.ResponseWriter, status int, v []byte)
}
