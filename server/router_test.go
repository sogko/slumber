package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/middlewares/renderer"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/test_helpers"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Router", func() {
	var s *server.Server
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var bodyJSON map[string]interface{}
	handleStub := func(ctx domain.IContext, version string) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			renderer.GetRendererCtx(ctx, req).Render(w, req, http.StatusOK, map[string]interface{}{
				"version": version,
			})
		}
	}
	handleAcl := func(req *http.Request, user domain.IUser) (bool, string) {
		return true, ""
	}

	r := renderer.New(&renderer.Options{IndentJSON: true}, renderer.JSON)

	Describe("Test API versioning", func() {

		var route domain.Route

		BeforeEach(func() {
			ctx := context.New()
			route = domain.Route{
				Name:           "Test",
				Method:         "GET",
				Pattern:        "/api/test",
				DefaultVersion: "0.2",
				RouteHandlers: domain.RouteHandlers{
					"0.1": handleStub(ctx, "0.1"),
					"0.2": handleStub(ctx, "0.2"),
					"0.3": handleStub(ctx, "0.3"),
				},
				ACLHandler: handleAcl,
			}
			routes := &domain.Routes{route}
			s = server.NewServer(&server.Config{
				Context: ctx,
			})
			router := server.NewRouter(ctx, nil)
			router.AddRoutes(routes)
			s.UseContextMiddleware(r)
			s.UseRouter(router)

			// record HTTP responses
			recorder = httptest.NewRecorder()
		})

		Context("when user does not specify API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})

		})

		Context("when user specify a valid Accept header (application/json) with valid API version", func() {

			It("should use specified API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.1")
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string("0.1")))
			})
		})

		Context("when user specify a valid Accept header (`vnd` tree + suffix case) with valid API version", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/vnd.api+json;version=0.10")
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header with API version that does not exists", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json;version=0.10")
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a valid Accept header but API version not specified", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "application/json")
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})

		Context("when user specify a empty Accept header", func() {

			It("should use default API version", func() {

				request, _ = http.NewRequest("GET", "/api/test", nil)
				request.Header.Set("Accept", "")
				s.ServeHTTP(recorder, request)
				bodyJSON = test_helpers.MapFromJSON(recorder.Body.Bytes())

				Expect(bodyJSON["version"]).To(Equal(string(route.DefaultVersion)))
			})
		})
	})

	Describe("AddRoutes()", func() {
		Context("Bad routes definition (undefined)", func() {
			It("should not panic", func() {
				Expect(func() {
					ctx := context.New()
					router := server.NewRouter(ctx, nil)
					router.AddRoutes(nil)
				}).ShouldNot(Panic())
			})

		})
		Context("Bad routes definition (missing default version handler)", func() {
			It("should panic", func() {
				Expect(func() {
					ctx := context.New()
					routes := &domain.Routes{
						domain.Route{
							Name:           "Test",
							Method:         "GET",
							Pattern:        "/api/test",
							DefaultVersion: "missinghandler",
							RouteHandlers: domain.RouteHandlers{
								"0.1": handleStub(ctx, "0.1"),
								"0.2": handleStub(ctx, "0.2"),
								"0.3": handleStub(ctx, "0.3"),
							},
							ACLHandler: handleAcl,
						},
					}
					router := server.NewRouter(ctx, nil)
					router.AddRoutes(routes)
				}).Should(Panic())
			})

		})
	})
	Describe("AddResources()", func() {
		Context("Valid IResource", func() {
			It("should not panic", func() {
				Expect(func() {
					ctx := context.New()
					router := server.NewRouter(ctx, nil)
					testResources := test_helpers.NewTestResource(ctx, r, &test_helpers.TestResourceOptions{})
					router.AddResources(testResources)
				}).ShouldNot(Panic())
			})

		})
		Context("Invalid IResource", func() {
			It("should panic", func() {
				Expect(func() {
					ctx := context.New()
					router := server.NewRouter(ctx, nil)
					testResources := test_helpers.NewTestResource(ctx, r, &test_helpers.TestResourceOptions{
						NilRoutes: true,
					})
					router.AddResources(testResources)
				}).Should(Panic())
			})

		})
	})

})
