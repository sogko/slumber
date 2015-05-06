package server_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Controller helpers", func() {
	var server *Server
	var route Route

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
	route = Route{
		Name:           "Test",
		Method:         "POST",
		Pattern:        "/api/test",
		DefaultVersion: "0.2",
		RouteHandlers: RouteHandlers{
			"0.2": func(w http.ResponseWriter, req *http.Request) {
				r := RendererCtx(req)
				var body TestBody
				err := DecodeJSONBodyHelper(w, req, r, &body)
				if err != nil {
					return
				}
				if body.A == "ERRORTEST" {
					RenderErrorResponseHelper(w, req, r, body.A)
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
		server = NewServer(&Config{
			Database: &DatabaseOptions{
				ServerName:   TestDatabaseServerName,
				DatabaseName: TestDatabaseName,
			},
			Renderer: &RendererOptions{},
			Routes:   &Routes{route},
		})

		// record HTTP responses
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
	})

	Describe("DecodeJSONBodyHelper()", func() {

		Context("when request body is a malformed JSON", func() {
			var response GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte("NOT A JSON")))
				server.ServeHTTP(recorder, request)
				utils.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeFalse())
			})

		})
		Context("when request body is a malformed JSON", func() {
			var response GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(`{"a": "OK"}`)))
				server.ServeHTTP(recorder, request)
				utils.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeTrue())
				Expect(response.Message).To(Equal("OK"))
			})

		})
	})

	Describe("RenderErrorResponseHelper()", func() {

		Context("when it is invoked", func() {
			var response GeneralResponse_v0
			BeforeEach(func() {
				request, _ = http.NewRequest("POST", "/api/test", bytes.NewReader([]byte(`{"a": "ERRORTEST"}`)))
				server.ServeHTTP(recorder, request)
				utils.DecodeResponseToType(recorder, &response)
			})

			It("returns error response", func() {
				Expect(response.Success).To(BeFalse())
				Expect(response.Message).To(Equal("ERRORTEST"))
			})

		})
	})
})
