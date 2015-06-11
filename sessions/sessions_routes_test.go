package sessions_test

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/repositories"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/sessions"
	th "github.com/sogko/slumber/test_helpers"
	"net/http"
	"net/http/httptest"
	"time"
)

const RequestAcceptHeader = "application/json;version=0.0"
const TestValidPassword = "PASSWORD1234"
const TestInvalidPassword = "PASSWORD1234WRONG"

var _ = Describe("Sessions API - /api/sessions; version=0.0", func() {

	var ts *th.TestServer
	var ts2 *th.TestServer
	var db domain.IDatabase
	var ta domain.ITokenAuthority
	var recorder *httptest.ResponseRecorder

	stubControllerHook := func(name string, doControllerHooksSuccess bool) func(w http.ResponseWriter, req *http.Request, ctx domain.IContext, payload interface{}) error {
		return func(w http.ResponseWriter, req *http.Request, ctx domain.IContext, payload interface{}) error {
			if payload == nil {
				return errors.New("Missing payload")
			}
			if !doControllerHooksSuccess {
				return errors.New("Expected hook to failed")
			}
			return nil
		}
	}

	BeforeEach(func() {

		// create test server
		ts = th.NewTestServer(&th.TestServerOptions{
			RequestAcceptHeader: RequestAcceptHeader,
			ServerConfig: &server.Config{
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
				Routes: &sessions.Routes,
				ACLMap: &sessions.ACL,
				ControllerHooks: &domain.ControllerHooksMap{
					PostCreateSessionHook: stubControllerHook("PostCreateSessionHook", true),
					PostDeleteSessionHook: stubControllerHook("PostDeleteSessionHook", true),
				},
			},
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		})

		// create another test server to test failed cases for controller hooks
		ts2 = th.NewTestServer(&th.TestServerOptions{
			RequestAcceptHeader: RequestAcceptHeader,
			ServerConfig: &server.Config{
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
				Routes: &sessions.Routes,
				ACLMap: &sessions.ACL,
				ControllerHooks: &domain.ControllerHooksMap{
					PostCreateSessionHook: stubControllerHook("PostCreateSessionHook", false),
					PostDeleteSessionHook: stubControllerHook("PostDeleteSessionHook", false),
				},
			},
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		})

		// create a separate db session so we can drop db later
		db = middlewares.NewMongoDB(&middlewares.MongoDBOptions{
			ServerName:   TestDatabaseServerName,
			DatabaseName: TestDatabaseName,
		})
		_ = db.NewSession()

		// setup token authority
		ta = middlewares.NewTokenAuthority(&middlewares.TokenAuthorityOptions{
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		})

		// record HTTP responses
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// drop database after each test
		db.DropDatabase()
	})
	Describe("GET /api/sessions", func() {
		Context("when token is still valid", func() {
			var response sessions.GetSessionResponse_v0
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

				ts.Request(recorder, "GET", "/api/sessions", nil, &response, &th.AuthOptions{Token: token})
			})
			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns response", func() {
				Expect(response.User.ID.Hex()).To(Equal(user.ID.Hex()))
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when request is non-authenticated", func() {
			var response sessions.GetSessionResponse_v0
			BeforeEach(func() {
				ts.Request(recorder, "GET", "/api/sessions", nil, &response, nil)
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

			Context("when PostCreateSessionHook returns OK", func() {
				var response sessions.CreateSessionResponse_v0
				BeforeEach(func() {
					ts.Request(recorder, "POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: user.Username,
						Password: TestValidPassword,
					}, &response, nil)

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

			Context("when PostCreateSessionHook failed", func() {
				var response sessions.GetSessionResponse_v0
				BeforeEach(func() {
					ts2.Request(recorder, "POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: user.Username,
						Password: TestValidPassword,
					}, &response, nil)
				})
				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
			})

		})

		Context("when user credentials is invalid", func() {
			Context("when password is wrong", func() {
				var response sessions.CreateSessionResponse_v0

				BeforeEach(func() {
					ts.Request(recorder, "POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: user.Username,
						Password: TestInvalidPassword,
					}, &response, nil)

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
					ts.Request(recorder, "POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: "",
						Password: TestInvalidPassword,
					}, &response, nil)
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
					ts.Request(recorder, "POST", "/api/sessions", sessions.CreateSessionRequest_v0{
						Username: "BOBBYSOX",
						Password: TestInvalidPassword,
					}, &response, nil)

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
				ts.Request(recorder, "POST", "/api/sessions", "BADJSON", &response, nil)
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
			Context("when PostDeleteSessionHook return OK", func() {
				var responseDelete sessions.DeleteSessionResponse_v0
				BeforeEach(func() {
					ts.Request(recorder, "DELETE", "/api/sessions", nil, &responseDelete, &th.AuthOptions{Token: token})
				})
				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns response", func() {
					Expect(responseDelete.Success).To(Equal(true))
				})
			})
			Context("when PostDeleteSessionHook failed", func() {
				var responseDelete sessions.DeleteSessionResponse_v0
				BeforeEach(func() {
					ts2.Request(recorder, "DELETE", "/api/sessions", nil, &responseDelete, &th.AuthOptions{Token: token})
				})
				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})

		Context("when token is does not have an identifier (`jti` claim)", func() {
			Context("when PostDeleteSessionHook returns OK", func() {
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
					ts.Request(recorder, "DELETE", "/api/sessions", nil, &responseDelete, &th.AuthOptions{Token: token})
				})
				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns response", func() {
					Expect(responseDelete.Success).To(Equal(true))
				})
			})
			Context("when PostDeleteSessionHook failed", func() {
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
					ts2.Request(recorder, "DELETE", "/api/sessions", nil, &responseDelete, &th.AuthOptions{Token: token})
				})
				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
			})
		})

		Context("when user try to login with previously revoked token", func() {
			var responseDelete sessions.DeleteSessionResponse_v0
			var responseDeleteSecondTime sessions.DeleteSessionResponse_v0
			var recorderSecond *httptest.ResponseRecorder
			BeforeEach(func() {
				ts.Request(recorder, "DELETE", "/api/sessions", nil, &responseDelete, &th.AuthOptions{Token: token})

				recorderSecond = httptest.NewRecorder()
				ts.Request(recorderSecond, "DELETE", "/api/sessions", nil, &responseDeleteSecondTime, &th.AuthOptions{Token: token})
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
