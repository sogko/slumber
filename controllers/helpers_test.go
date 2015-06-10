package controllers_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/controllers"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/test_helpers"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Controller helpers", func() {
	var s *server.Server
	var route domain.Route

	var request *http.Request
	var recorder *httptest.ResponseRecorder

	type TestBody struct {
		A string `json:"a"`
	}
	type TestResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	// define test route
	route = domain.Route{
		Name:           "Test",
		Method:         "POST",
		Pattern:        "/api/test",
		DefaultVersion: "0.2",
		RouteHandlers: domain.RouteHandlers{
			"0.2": func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
				r := ctx.GetRendererCtx(req)
				var body TestBody
				err := controllers.DecodeJSONBodyHelper(w, req, r, &body)
				if err != nil {
					return
				}
				if body.A == "ERRORTEST" {
					controllers.RenderErrorResponseHelper(w, req, r, body.A)
					return
				}

				r.JSON(w, http.StatusOK, TestResponse{
					Success: true,
					Message: body.A,
				})

			},
		},
	}

	BeforeEach(func() {
		s = server.NewServer(&server.Config{
			Database: &middlewares.MongoDBOptions{
				ServerName:   TestDatabaseServerName,
				DatabaseName: TestDatabaseName,
			},
			Renderer: &middlewares.RendererOptions{},
			Routes:   &domain.Routes{route},
			ACLMap: &domain.ACLMap{
				"Test": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
					return true, ""
				},
			},
		}).SetupRoutes()

		// record HTTP responses
		recorder = httptest.NewRecorder()
	})

	Describe("DecodeJSONBodyHelper()", func() {

		Context("when request body is a malformed JSON", func() {
			var response controllers.GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte("NOT A JSON")))
				s.ServeHTTP(recorder, request)
				test_helpers.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeFalse())
			})

		})
		Context("when request body is a malformed JSON", func() {
			var response controllers.GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(`{"a": "OK"}`)))
				s.ServeHTTP(recorder, request)
				test_helpers.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeTrue())
				Expect(response.Message).To(Equal("OK"))
			})

		})
	})

	Describe("RenderErrorResponseHelper()", func() {

		Context("when it is invoked", func() {
			var response controllers.GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(`{"a": "ERRORTEST"}`)))
				s.ServeHTTP(recorder, request)
				test_helpers.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeFalse())
				Expect(response.Message).To(Equal("ERRORTEST"))
			})

		})
	})
})
