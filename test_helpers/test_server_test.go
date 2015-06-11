package test_helpers_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/test_helpers"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
)

type TestRequestBody struct {
	Value string
}
type TestResponseBody struct {
	Result string
	Value  string
}

var _ = Describe("Test Server", func() {

	var ts *test_helpers.TestServer

	Describe("Request()", func() {

		BeforeEach(func() {
			// configure test server
			serverConfig := &server.Config{
				Database: &middlewares.MongoDBOptions{
					ServerName:   TestDatabaseServerName,
					DatabaseName: TestDatabaseName,
				},
				Renderer: &middlewares.RendererOptions{
					IndentJSON: true,
				},
				TokenAuthority: &middlewares.TokenAuthorityOptions{
					PrivateSigningKey: privateSigningKey,
					PublicSigningKey:  publicSigningKey,
				},
				Routes: &domain.Routes{
					domain.Route{
						Name:           "TestGetRoute",
						Method:         "GET",
						Pattern:        "/api/test",
						DefaultVersion: "0.0",
						RouteHandlers: domain.RouteHandlers{
							"0.0": func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
								r := ctx.GetRendererCtx(req)
								r.JSON(w, http.StatusOK, TestResponseBody{
									Result: "OK",
								})
							},
						},
					},
					domain.Route{
						Name:           "TestPostRoute",
						Method:         "POST",
						Pattern:        "/api/test",
						DefaultVersion: "0.0",
						RouteHandlers: domain.RouteHandlers{
							"0.0": func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
								r := ctx.GetRendererCtx(req)

								var body TestRequestBody
								decoder := json.NewDecoder(req.Body)
								err := decoder.Decode(&body)
								if err != nil {
									r.JSON(w, http.StatusBadRequest, TestResponseBody{
										Result: "NOT_OK",
									})
								}
								r.JSON(w, http.StatusOK, TestResponseBody{
									Result: "OK",
									Value:  body.Value,
								})
							},
						},
					},
				},
				ACLMap: &domain.ACLMap{
					"TestGetRoute": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
						return true, ""
					},
					"TestPostRoute": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
						return true, ""
					},
				},
			}

			// create test server
			ts = test_helpers.NewTestServer(&test_helpers.TestServerOptions{
				RequestAcceptHeader: "application/json;version=0.0",
				ServerConfig:        serverConfig,
				PrivateSigningKey:   privateSigningKey,
				PublicSigningKey:    publicSigningKey,
			})

		})

		Context("Basic request", func() {
			It("returns status code of StatusOK (200)", func() {
				var response TestResponseBody
				recorder := httptest.NewRecorder()

				ts.Request(recorder, "GET", "/api/test", nil, &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.Result).To(Equal("OK"))
			})
		})

		Context("Non-empty JSON valid body", func() {
			It("returns status code of StatusOK (200)", func() {
				var response TestResponseBody
				recorder := httptest.NewRecorder()

				ts.Request(recorder, "POST", "/api/test", TestRequestBody{
					Value: "string",
				}, &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusOK))
				Expect(response.Result).To(Equal("OK"))
				Expect(response.Value).To(Equal("string"))
			})
		})
		Context("Non-empty JSON invalid body", func() {
			It("returns status code of StatusBadRequest (400)", func() {
				var response TestResponseBody
				recorder := httptest.NewRecorder()

				ts.Request(recorder, "POST", "/api/test", "INVALID", &response, nil)

				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				Expect(response.Result).To(Equal("NOT_OK"))
			})
		})

		Context("AuthOptions.Token", func() {
			It("returns status code of StatusUnauthorized (401)", func() {
				var response TestResponseBody
				recorder := httptest.NewRecorder()
				ts.Request(recorder, "GET", "/api/test", nil, &response, &test_helpers.AuthOptions{
					Token: "invalidrandomtokenshould401",
				})
				Expect(recorder.Code).To(Equal(http.StatusUnauthorized))
			})
		})
		Context("AuthOptions.APIUser", func() {
			It("returns status code of StatusOK (400)", func() {

				// create fake user
				// since routes does not need authorization to access
				user := domain.User{
					ID:       bson.NewObjectId(),
					Username: "admin",
					Status:   "active",
					Roles: domain.Roles{
						domain.RoleAdmin,
						domain.RoleUser,
					},
				}

				var response TestResponseBody
				recorder := httptest.NewRecorder()
				ts.Request(recorder, "GET", "/api/test", nil, &response, &test_helpers.AuthOptions{
					APIUser: &user,
				})
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
		})
	})

})
