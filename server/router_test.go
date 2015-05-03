package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
	"net/http/httptest"
)

func handleStub(version string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r := RendererCtx(req)
		r.JSON(w, http.StatusOK, map[string]interface{}{
			"version": version,
		})
	}
}

var _ = Describe("Router", func() {
	var server *Server
	var session *DatabaseSession
	var renderer *Renderer
	var route Route
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var bodyJSON map[string]interface{}

	// set up server with test components
	session = NewSession(DatabaseOptions{
		ServerName:   TestDatabaseServerName,
		DatabaseName: TestDatabaseName,
	})
	renderer = NewRenderer(RendererOptions{})

	route = Route{
		Name:           "Test",
		Method:         "GET",
		Pattern:        "/api/test",
		DefaultVersion: "0.2",
		RouteHandlers: RouteHandlers{
			"0.1": handleStub("0.1"),
			"0.2": handleStub("0.2"),
			"0.3": handleStub("0.3"),
		},
	}

	Describe("Test API versioning", func() {

		BeforeEach(func() {
			routes := &Routes{route}
			server = NewServer(&Components{
				DatabaseSession: session,
				Renderer:        renderer,
				Routes:          routes,
			})

			// record HTTP responses
			recorder = httptest.NewRecorder()
		})

		Context("when user does not specify API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})

		})

		Context("when user specify a valid Accept header (application/json) with valid API version", func() {

			It("should use specified API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.1")
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string("0.1")))
			})
		})

		Context("when user specify a valid Accept header (`vnd` tree + suffix case) with valid API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/vnd.api+json;version=0.10")
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header with API version that does not exists", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.10")
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header but API version not specified", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json")
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a empty Accept header", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "")
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})
	})

	Describe("Bad routes definition (undefined)", func() {

		It("should panic", func() {
			Expect(func() {
				server = NewServer(&Components{
					DatabaseSession: session,
					Renderer:        renderer,
				})

			}).Should(Panic())
		})

	})
	Describe("Bad routes definition (missing default version handler)", func() {

		It("should panic", func() {
			Expect(func() {
				routes := &Routes{
					Route{"Test", "GET", "/api/test", "missingDefaultVersion", RouteHandlers{
						"0.1": handleStub("0.1"),
					}},
				}
				server = NewServer(&Components{
					DatabaseSession: session,
					Renderer:        renderer,
					Routes:          routes,
				})

			}).Should(Panic())
		})

	})

})
