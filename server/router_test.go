package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/libs"
	"github.com/sogko/golang-rest-api-server-example/middlewares"
	"github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
	"net/http/httptest"
)

func handleStub(version string) domain.ContextHandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
		r := ctx.GetRendererCtx(req)
		r.JSON(w, http.StatusOK, map[string]interface{}{
			"version": version,
		})
	}
}

var _ = Describe("Router", func() {
	var s *server.Server
	var route domain.Route
	var aclMap domain.ACLMap
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var bodyJSON map[string]interface{}

	var dbOptions = middlewares.MongoDBOptions{
		ServerName:   TestDatabaseServerName,
		DatabaseName: TestDatabaseName,
	}

	var renderOptions = middlewares.RendererOptions{}

	// define test route
	route = domain.Route{
		Name:           "Test",
		Method:         "GET",
		Pattern:        "/api/test",
		DefaultVersion: "0.2",
		RouteHandlers: domain.RouteHandlers{
			"0.1": handleStub("0.1"),
			"0.2": handleStub("0.2"),
			"0.3": handleStub("0.3"),
		},
	}

	aclMap = domain.ACLMap{
		"Test": func(user *domain.User, req *http.Request, ctx domain.IContext) bool {
			return (true)
		},
	}

	Describe("Test API versioning", func() {

		BeforeEach(func() {
			routes := &domain.Routes{route}
			s = server.NewServer(&server.Config{
				Database:       &dbOptions,
				Renderer:       &renderOptions,
				Routes:         routes,
				TokenAuthority: &middlewares.TokenAuthorityOptions{},
				ACLMap:         &aclMap,
			})

			// record HTTP responses
			recorder = httptest.NewRecorder()
		})

		Context("when user does not specify API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})

		})

		Context("when user specify a valid Accept header (application/json) with valid API version", func() {

			It("should use specified API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.1")
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string("0.1")))
			})
		})

		Context("when user specify a valid Accept header (`vnd` tree + suffix case) with valid API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/vnd.api+json;version=0.10")
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header with API version that does not exists", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.10")
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header but API version not specified", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json")
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a empty Accept header", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "")
				s.ServeHTTP(recorder, request)
				bodyJSON = libs.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})
	})

	Describe("Bad routes definition (undefined)", func() {

		It("should panic", func() {
			Expect(func() {
				s = server.NewServer(&server.Config{
					Database: &dbOptions,
					Renderer: &renderOptions,
				})

			}).Should(Panic())
		})

	})
	Describe("Bad routes definition (missing default version handler)", func() {

		It("should panic", func() {
			Expect(func() {
				routes := &domain.Routes{
					domain.Route{"Test", "GET", "/api/test", "missingDefaultVersion", domain.RouteHandlers{
						"0.1": handleStub("0.1"),
					}},
				}
				s = server.NewServer(&server.Config{
					Database: &dbOptions,
					Renderer: &renderOptions,
					Routes:   routes,
				})

			}).Should(Panic())
		})

	})

})
