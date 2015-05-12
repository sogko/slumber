package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"net/http"
)

type GeneralResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

type ErrorResponse_v0 struct {
	Message string `json:"message,omitempty"`
	Success bool   `json:"success"`
}

// DecodeJSONBodyHelper is a helper function to decode JSON request body
func DecodeJSONBodyHelper(w http.ResponseWriter, req *http.Request, r domain.IRenderer, target interface{}) error {
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(target)
	if err != nil {
		RenderErrorResponseHelper(w, req, r, fmt.Sprintf("Request body parse error: %v", err.Error()))
		return err
	}
	return nil
}

// RenderErrorResponseHelper is a helper function to render consistent error message
func RenderErrorResponseHelper(w http.ResponseWriter, req *http.Request, r domain.IRenderer, message string) {
	r.JSON(w, http.StatusBadRequest, ErrorResponse_v0{
		Message: message,
		Success: false,
	})
}
