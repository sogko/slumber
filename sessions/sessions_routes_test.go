package sessions_test

import (
	"bytes"
	"encoding/json"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example"
	. "github.com/sogko/golang-rest-api-server-example/server"
	. "github.com/sogko/golang-rest-api-server-example/sessions"
	"github.com/sogko/golang-rest-api-server-example/users"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"net/http"
	"net/http/httptest"
)

const RequestAcceptHeader = "application/json;version=0.0"
const TestValidPassword = "PASSWORD1234"

var _ = Describe("Sessions API - /api/sessions; version=0.0", func() {

	var server *Server
	var session *DatabaseSession
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	// sendRequestHelper helps
	// - creates an HTTP request
	// - set the right API version
	// - serves request
	// - decode the response into the desired interface
	sendRequestHelper := func(method string, urlStr string, body interface{}, targetResponse interface{}) {

		// request for version 0.0
		if body != nil {
			jsonBytes, _ := json.Marshal(body)
			request, _ = http.NewRequest(method, urlStr, bytes.NewReader(jsonBytes))
		} else {
			request, _ = http.NewRequest(method, urlStr, nil)
		}

		// set API version through accept header
		request.Header.Set("Accept", RequestAcceptHeader)

		// serve request
		server.ServeHTTP(recorder, request)
		utils.DecodeResponseToType(recorder, &targetResponse)
	}

	BeforeEach(func() {

		dbOptions := DatabaseOptions{
			ServerName:   TestDatabaseServerName,
			DatabaseName: TestDatabaseName,
		}

		// create a separate db session so we can drop db later
		session = NewSession(dbOptions)

		// init server
		server = NewServer(&Config{
			Database: &dbOptions,
			Renderer: &RendererOptions{
				IndentJSON: true,
			},
			Routes: GetRoutes(),
		})

		// record HTTP responses
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// drop database after each test
		session.DB(session.DatabaseName).DropDatabase()
	})

	Describe("POST /api/sessions", func() {

		var user *users.User
		BeforeEach(func() {
			// insert
			user = gory.Build("user").(*users.User)
			user.SetPassword(TestValidPassword)
			collection := session.DB(TestDatabaseName).C(users.UsersCollection)
			collection.Insert(user)
		})
		Context("test", func() {
			var response CreateResponse_v0

			BeforeEach(func() {

				sendRequestHelper("POST", "/api/sessions", CreateRequest_v0{
					Username: user.Username,
					Password: TestValidPassword, // test password
				}, &response)

			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns token", func() {

				token, claims, _ := utils.ParseAndVerifyTokenString(response.Token)

				Expect(response.Success).To(Equal(true))

				Expect(token.Valid).To(Equal(true))
				Expect(claims.ID).To(Equal(user.ID.Hex()))
				Expect(claims.Status).To(Equal(user.Status))
				Expect(len(claims.Roles)).To(Equal(len(user.Roles)))
			})
		})

	})

})
