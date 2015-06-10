package sessions_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/repositories"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/sessions"
	"github.com/sogko/slumber/test_helpers"
	"net/http"
	"net/http/httptest"
	"time"
)

const RequestAcceptHeader = "application/json;version=0.0"
const TestValidPassword = "PASSWORD1234"
const TestInvalidPassword = "PASSWORD1234WRONG"

var _ = Describe("Sessions API - /api/sessions; version=0.0", func() {

	var s *server.Server
	var db domain.IDatabase
	var ta domain.ITokenAuthority
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	// sendRequestHelper helps
	// - creates an HTTP request
	// - set the right API version
	// - serves request
	// - decode the response into the desired interface
	sendRequestHelper := func(method string, urlStr string, body interface{}, apiUser *domain.User, targetResponse interface{}) {

		// request for version 0.0
		if body != nil {
			jsonBytes, _ := json.Marshal(body)
			request, _ = http.NewRequest(method, urlStr, bytes.NewReader(jsonBytes))
		} else {
			request, _ = http.NewRequest(method, urlStr, nil)
		}

		// set API version through accept header
		request.Header.Set("Accept", RequestAcceptHeader)

		if apiUser != nil {
			// set Authorization header
			// TODO: create utility
			var rolesString []string
			for _, role := range apiUser.Roles {
				rolesString = append(rolesString, string(role))
			}
			token, _ := ta.CreateNewSessionToken(&domain.TokenClaims{
				UserID:   apiUser.ID.Hex(),
				Username: apiUser.Username,
				Status:   apiUser.Status,
				Roles:    rolesString,
			})
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		}

		// serve request
		s.ServeHTTP(recorder, request)
		test_helpers.DecodeResponseToType(recorder, &targetResponse)
	}
	sendRequestWithTokenHelper := func(recorder *httptest.ResponseRecorder, method string, urlStr string, body interface{}, token string, targetResponse interface{}) {

		// request for version 0.0
		if body != nil {
			jsonBytes, _ := json.Marshal(body)
			request, _ = http.NewRequest(method, urlStr, bytes.NewReader(jsonBytes))
		} else {
			request, _ = http.NewRequest(method, urlStr, nil)
		}
		// set API version through accept header
		request.Header.Set("Accept", RequestAcceptHeader)
		if token != "" {
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", token))
		}

		// serve request
		s.ServeHTTP(recorder, request)
		test_helpers.DecodeResponseToType(recorder, &targetResponse)
	}
	BeforeEach(func() {

		dbOptions := middlewares.MongoDBOptions{
			ServerName:   TestDatabaseServerName,
			DatabaseName: TestDatabaseName,
		}
		// create a separate db session so we can drop db later
		db = middlewares.NewMongoDB(&dbOptions)
		_ = db.NewSession()

		// init server
		s = server.NewServer(&server.Config{
			Database: &dbOptions,
			Renderer: &middlewares.RendererOptions{
				IndentJSON: true,
			},
			TokenAuthority: &middlewares.TokenAuthorityOptions{
				PrivateSigningKey: privateSigningKey,
				PublicSigningKey:  publicSigningKey,
			},
			Routes: &sessions.Routes,
			ACLMap: &sessions.ACL,
		}).SetupRoutes()

		// record HTTP responses
		recorder = httptest.NewRecorder()

		// setup token authority
		ta = middlewares.NewTokenAuthority(&middlewares.TokenAuthorityOptions{
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		})
	})

	AfterEach(func() {
		// drop database after each test
		db.DropDatabase()
	})
	Describe("GET /api/sessions", func() {
		var token string

		var user *domain.User
		BeforeEach(func() {
			// insert a user
			user = gory.Build("user").(*domain.User)
			user.SetPassword(TestValidPassword)
			db.Insert(repositories.UsersCollection, user)

			// create a session token
			token, _ = ta.CreateNewSessionToken(&domain.TokenClaims{
				UserID:   user.ID.Hex(),
				Username: user.Username,
				Status:   "active",
				Roles:    []string{"admin"},
			})

		})

		Context("when token is still valid", func() {
			var response sessions.GetSessionResponse_v0
			BeforeEach(func() {
				sendRequestWithTokenHelper(recorder, "GET", "/api/sessions", nil, token, &response)
			})
			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns response", func() {
				Expect(response.User.ID).To(Equal(user.ID))
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when request is non-authenticated", func() {
			var response sessions.GetSessionResponse_v0
			BeforeEach(func() {
				sendRequestWithTokenHelper(recorder, "GET", "/api/sessions", nil, "", &response)
			})
			It("returns status code of StatusForbidden (403)", func() {
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
			It("returns response", func() {
				Expect(response.Success).To(Equal(false))
			})
		})
	})
	Describe("POST /api/sessions", func() {

		var user *domain.User
		BeforeEach(func() {
			// insert a user
			user = gory.Build("user").(*domain.User)
			user.SetPassword(TestValidPassword)
			db.Insert(repositories.UsersCollection, user)
		})
		Context("when user credentials is valid", func() {
			var response sessions.CreateSessionResponse_v0

			BeforeEach(func() {
				sendRequestHelper("POST", "/api/sessions", sessions.CreateSessionRequest_v0{
					Username: user.Username,
					Password: TestValidPassword,
				}, nil, &response)

			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns a token valid", func() {

				token, claims, _ := ta.VerifyTokenString(response.Token)

				Expect(response.Success).To(Equal(true))

				Expect(token.Valid).To(Equal(true))
				Expect(claims.UserID).To(Equal(user.ID.Hex()))
				Expect(claims.Username).To(Equal(user.Username))
				Expect(claims.Status).To(Equal(user.Status))
				Expect(len(claims.Roles)).To(Equal(len(user.Roles)))
			})
		})

		Context("when user credentials is invalid", func() {
			Context("when password is wrong", func() {
				var response sessions.CreateSessionResponse_v0

				BeforeEach(func() {
					sendRequestHelper("POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: user.Username,
						Password: TestInvalidPassword,
					}, nil, &response)

				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("does not return a token valid", func() {
					Expect(response.Success).To(Equal(false))
				})
			})
			Context("when username is empty", func() {
				var response sessions.CreateSessionResponse_v0

				BeforeEach(func() {
					sendRequestHelper("POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: "",
						Password: TestInvalidPassword,
					}, nil, &response)

				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("does not return a token valid", func() {
					Expect(response.Success).To(Equal(false))
				})
			})
			Context("when username does not exists", func() {
				var response sessions.CreateSessionResponse_v0

				BeforeEach(func() {
					sendRequestHelper("POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: "BOBBYSOX",
						Password: TestInvalidPassword,
					}, nil, &response)

				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("does not return a token valid", func() {
					Expect(response.Success).To(Equal(false))
				})
			})
		})

		Context("when POSTing a malformed JSON body", func() {

			var response sessions.CreateSessionResponse_v0

			BeforeEach(func() {
				sendRequestHelper("POST", "/api/sessions", "BADJSON", nil, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("does not return a token valid", func() {
				Expect(response.Success).To(Equal(false))
			})
		})

	})
	Describe("DELETE /api/sessions", func() {

		var token string

		var user *domain.User
		BeforeEach(func() {
			// insert a user
			user = gory.Build("user").(*domain.User)
			user.SetPassword(TestValidPassword)
			db.Insert(repositories.UsersCollection, user)

			// create a session token
			token, _ = ta.CreateNewSessionToken(&domain.TokenClaims{
				UserID:   user.ID.Hex(),
				Username: user.Username,
				Status:   "active",
				Roles:    []string{"admin"},
			})

		})

		Context("when token is still valid and has not been revoked previously", func() {
			var responseDelete sessions.DeleteSessionResponse_v0
			BeforeEach(func() {
				sendRequestWithTokenHelper(recorder, "DELETE", "/api/sessions", nil, token, &responseDelete)
			})
			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns response", func() {
				Expect(responseDelete.Success).To(Equal(true))
			})
		})

		Context("when token is does not have an identifier (`jti` claim)", func() {
			var responseDelete sessions.DeleteSessionResponse_v0
			BeforeEach(func() {

				// insert a user
				user = gory.Build("user").(*domain.User)
				user.SetPassword(TestValidPassword)
				db.Insert(repositories.UsersCollection, user)

				// manually generate jwt so that we can set an invalid JTI claim
				// possible from a malicious attacker
				tokenObj := jwt.New(jwt.SigningMethodRS512)
				tokenObj.Claims = map[string]interface{}{
					"userId": user.ID.Hex(),
					"status": "active",
					"roles":  []string{"admin"},
					"exp":    time.Now().Add(time.Hour * 72).Format(time.RFC3339), // 3 days
					"iat":    time.Now().Format(time.RFC3339),
					"jti":    "INVALIDJTI",
				}
				token, _ := tokenObj.SignedString(privateSigningKey)
				sendRequestWithTokenHelper(recorder, "DELETE", "/api/sessions", nil, token, &responseDelete)
			})
			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns response", func() {
				Expect(responseDelete.Success).To(Equal(true))
			})
		})

		Context("when user try to login with previously revoked token", func() {
			var responseDelete sessions.DeleteSessionResponse_v0
			var responseDeleteSecondTime sessions.DeleteSessionResponse_v0
			var recorderSecond *httptest.ResponseRecorder
			BeforeEach(func() {
				sendRequestWithTokenHelper(recorder, "DELETE", "/api/sessions", nil, token, &responseDelete)

				recorderSecond = httptest.NewRecorder()
				sendRequestWithTokenHelper(recorderSecond, "DELETE", "/api/sessions", nil, token, &responseDeleteSecondTime)
			})

			It("returns status code of StatusUnauthorized (401)", func() {
				Expect(recorderSecond.Code).To(Equal(http.StatusUnauthorized))
			})
			It("returns failed response", func() {
				Expect(responseDeleteSecondTime.Success).To(Equal(false))
			})
		})
	})
})
