package test_helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/sogko/slumber/domain"
)

type TestRequestBody struct {
	Value string
}
type TestResponseBody struct {
	Result string
	Value  string
}

type TestResourceOptions struct {
	NilRoutes bool
}

// TestResource implements IResource
func NewTestResource(ctx domain.IContext, r domain.IRenderer, options *TestResourceOptions) *TestResource {
	return &TestResource{ctx, r, options}
}

type TestResource struct {
	ctx      domain.IContext
	Renderer domain.IRenderer
	Options  *TestResourceOptions
}

func (resource *TestResource) Context() domain.IContext {
	return resource.ctx
}
func (resource *TestResource) Routes() *domain.Routes {
	if resource.Options.NilRoutes == true {
		return nil
	}
	return &domain.Routes{
		domain.Route{
			Name:           "TestGetRoute",
			Method:         "GET",
			Pattern:        "/api/test",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandleGetRoute,
			},
			ACLHandler: resource.HandleAllACL,
		},
		domain.Route{
			Name:           "TestPostRoute",
			Method:         "POST",
			Pattern:        "/api/test",
			DefaultVersion: "0.0",
			RouteHandlers: domain.RouteHandlers{
				"0.0": resource.HandlePostRoute,
			},
			ACLHandler: resource.HandleAllACL,
		},
	}
}
func (resource *TestResource) Render(w http.ResponseWriter, req *http.Request, status int, v interface{}) {
	resource.Renderer.JSON(w, status, v)
}
func (resource *TestResource) HandleAllACL(req *http.Request, user domain.IUser) (bool, string) {
	return true, ""
}
func (resource *TestResource) HandleGetRoute(w http.ResponseWriter, req *http.Request) {
	resource.Render(w, req, http.StatusOK, TestResponseBody{
		Result: "OK",
	})
}
func (resource *TestResource) HandlePostRoute(w http.ResponseWriter, req *http.Request) {
	var body TestRequestBody
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&body)
	if err != nil {
		resource.Render(w, req, http.StatusBadRequest, TestResponseBody{
			Result: "NOT_OK",
		})
	}
	resource.Render(w, req, http.StatusOK, TestResponseBody{
		Result: "OK",
		Value:  body.Value,
	})
}

func NewTestContextMiddleware() *TestContextMiddleware {
	return &TestContextMiddleware{}
}

type TestContextMiddleware struct {
}

func (middleware *TestContextMiddleware) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
	next(w, req)
}

func NewTestMiddleware() *TestMiddleware {
	return &TestMiddleware{}
}

type TestMiddleware struct {
}

func (middleware *TestMiddleware) Handler(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	next(w, req)
}

// MapFromJSON is a test helper function that decodes recorded response body to
// a specific struct type
// Note: this functions panics on error. For test usage only, not for production.
func MapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(fmt.Sprintf("mapFromJSON(): Not a valid JSON body\n%v", string(data)))
	}
	return result.(map[string]interface{})
}

// DecodeResponseToType is a test helper function that decodes recorded response body to
// a specific struct type
// Note: this functions panics on error. For test usage only, not for production.
func DecodeResponseToType(recorder *httptest.ResponseRecorder, target interface{}) error {
	// clone request body reader so that we can have a nicer error message
	bodyString := ""
	if b, err := ioutil.ReadAll(recorder.Body); err == nil {
		bodyString = string(b)
	}
	readerClone := strings.NewReader(bodyString)

	decoder := json.NewDecoder(readerClone)
	err := decoder.Decode(target)
	if err != nil {
		log.Println(fmt.Sprintf("DecodeResponseToType(): %v \n%v", err.Error(), bodyString))
	}
	return err
}
