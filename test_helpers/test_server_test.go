package test_helpers_test

import (
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber-sessions"
	"github.com/sogko/slumber-users"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/middlewares/renderer"
	"github.com/sogko/slumber/test_helpers"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Test Server", func() {

	var ts *test_helpers.TestServer

	Describe("Request()", func() {

		var sessionsResource *sessions.Resource = nil
		var usersResource *users.Resource = nil
		BeforeEach(func() {

			// configure test server
			ctx := context.New()

			renderer := renderer.New(&renderer.Options{
				IndentJSON: true,
			}, renderer.JSON)

			testMiddleware := test_helpers.NewTestMiddleware()
			testContextMiddleware := test_helpers.NewTestContextMiddleware()
			testResource := test_helpers.NewTestResource(ctx, renderer, &test_helpers.TestResourceOptions{})

			// create test server
			ts = test_helpers.NewTestServer(&test_helpers.TestServerOptions{
				RequestAcceptHeader: "application/json;version=0.0",
				PrivateSigningKey:   privateSigningKey,
				PublicSigningKey:    publicSigningKey,
				Resources:           []domain.IResource{testResource},
				Middlewares:         []interface{}{testMiddleware, testContextMiddleware, nil},
			})

			usersResource = users.NewResource(ctx, &users.Options{
				Database: ts.Database,
				Renderer: ts.Renderer,
			})

			sessionsResource = sessions.NewResource(ctx, &sessions.Options{
				PrivateSigningKey:     privateSigningKey,
				PublicSigningKey:      publicSigningKey,
				Database:              ts.Database,
				Renderer:              ts.Renderer,
				UserRepositoryFactory: usersResource.UserRepositoryFactory,
			})
			ts.AddResources(sessionsResource)

		})

		Context("Basic request", func() {
			It("returns status code of StatusOK (200)", func() {
				var response test_helpers.TestResponseBody
				ts.Run()
				recorder := ts.Request("GET", "/api/test", nil, &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.Result).To(Equal("OK"))
			})
		})

		Context("Non-empty JSON valid body", func() {
			It("returns status code of StatusOK (200)", func() {
				var response test_helpers.TestResponseBody
				ts.Run()
				recorder := ts.Request("POST", "/api/test", test_helpers.TestRequestBody{
					Value: "string",
				}, &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.Result).To(Equal("OK"))
				Expect(response.Value).To(Equal("string"))
			})
		})
		Context("Non-empty JSON invalid body", func() {
			It("returns status code of StatusBadRequest (400)", func() {
				var response test_helpers.TestResponseBody
				ts.Run()
				recorder := ts.Request("POST", "/api/test", "INVALID", &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(response.Result).To(Equal("NOT_OK"))
			})
		})

		Context("AuthOptions.Token", func() {
			Context("without sessions.Authenticator enabled", func() {
				It("returns status code of StatusUnauthorized (401)", func() {
					var response test_helpers.TestResponseBody
					ts.Run()
					recorder := ts.Request("GET", "/api/test", nil, &response, &test_helpers.AuthOptions{
						Token: "invalidrandomtokenshould401",
					})
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
			})
			Context("with sessions.Authenticator enabled", func() {
				It("returns status code of StatusUnauthorized (401)", func() {
					var response test_helpers.TestResponseBody

					// add sessions authenticator middleware
					ts.AddMiddlewares(sessionsResource.NewAuthenticator())
					ts.Run()

					recorder := ts.Request("GET", "/api/test", nil, &response, &test_helpers.AuthOptions{
						Token: "invalidrandomtokenshould401",
					})
					Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
				})
			})
		})
		Context("AuthOptions.APIUser", func() {
			It("returns status code of StatusOK (400)", func() {

				// create fake user
				// since routes does not need authorization to access
				user := users.User{
					ID:       bson.NewObjectId(),
					Username: "admin",
					Status:   "active",
					Roles: users.Roles{
						users.RoleAdmin,
						users.RoleUser,
					},
				}

				var response test_helpers.TestResponseBody
				ts.Run()
				recorder := ts.Request("GET", "/api/test", nil, &response, &test_helpers.AuthOptions{
					APIUser: &user,
				})
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})

})
