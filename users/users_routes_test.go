package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/libs"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/repositories"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/users"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"time"
)

const RequestAcceptHeader = "application/json;version=0.0"

var _ = Describe("Users API - /api/users; version=0.0", func() {

	var s *server.Server
	var db domain.IDatabase
	var ta domain.ITokenAuthority
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	// sendRequestHelper helps
	// - creates an HTTP request
	// - set the right API version
	// - serves request
	// - decode the response into the desired interfac
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
		libs.DecodeResponseToType(recorder, &targetResponse)
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
			Routes: &users.UsersAPIRoutes,
			ACLMap: &users.UsersAPIACL,
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

	Describe("GET /api/users", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})

			Context("when no users exist", func() {

				var response users.ListUsersResponse_v0

				BeforeEach(func() {
					sendRequestHelper("GET", "/api/users", nil, adminUser, &response)
				})

				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns zero users", func() {
					Expect(len(response.Users)).To(Equal(1))
					Expect(response.Success).To(Equal(true))
				})
			})

			Context("when two users exist", func() {

				var response users.ListUsersResponse_v0

				BeforeEach(func() {
					// insert two valid users into database
					db.Insert(repositories.UsersCollection, gory.Build("user"))
					db.Insert(repositories.UsersCollection, gory.Build("user"))

					sendRequestHelper("GET", "/api/users", nil, adminUser, &response)
				})

				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns two users", func() {
					Expect(len(response.Users)).To(Equal(3))
					Expect(response.Success).To(Equal(true))
				})

			})
		})

		Context("when API user is anonymous", func() {

			var response users.ListUsersResponse_v0
			It("returns status code of StatusOK (200)", func() {
				sendRequestHelper("GET", "/api/users", nil, nil, &response)
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("when API user is an inactive user", func() {

			var inactiveAdminUser *domain.User

			BeforeEach(func() {
				inactiveAdminUser = gory.Build("inactiveAdminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, inactiveAdminUser)
			})
			var response users.ListUsersResponse_v0
			It("returns status code of StatusOK (200)", func() {
				sendRequestHelper("GET", "/api/users", nil, inactiveAdminUser, &response)
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
		})
	})

	Describe("POST /api/users", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})

			Context("when adding one valid users", func() {

				var newUser *domain.NewUser
				var response users.CreateUserResponse_v0

				BeforeEach(func() {

					newUser = gory.Build("newUser").(*domain.NewUser)

					sendRequestHelper("POST", "/api/users", users.CreateUserRequest_v0{
						User: *newUser,
					}, adminUser, &response)

				})

				It("returns status code of StatusCreated (201)", func() {
					Expect(recorder.Code).To(Equal(http.StatusCreated))
				})
				It("returns newly-created user", func() {
					Expect(response.User.Email).To(Equal(newUser.Email))
					Expect(response.User.Roles).To(BeNil())
					Expect(response.User.Status).To(Equal(domain.StatusPending))
					Expect(response.Success).To(Equal(true))
				})
			})

			Context("when trying to specify Roles and Status", func() {

				var newUser *domain.NewUser
				var response users.CreateUserResponse_v0

				BeforeEach(func() {

					newUser = gory.Build("newUser").(*domain.NewUser)
					sendRequestHelper("POST", "/api/users", users.CreateUserRequest_v0{
						User: *newUser,
					}, adminUser, &response)

				})

				It("returns status code of StatusCreated (201)", func() {
					Expect(recorder.Code).To(Equal(http.StatusCreated))
				})
				It("returns newly-created user", func() {
					Expect(response.User.Email).To(Equal(newUser.Email))
					Expect(response.User.Roles).To(BeNil())
					Expect(response.User.Status).To(Equal(domain.StatusPending))
					Expect(response.Success).To(Equal(true))
				})
			})

			Context("when POSTing a malformed JSON body", func() {

				var response users.CreateUserResponse_v0

				BeforeEach(func() {
					sendRequestHelper("POST", "/api/users", "BADJSON", adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns nil user", func() {
					Expect(response.User).To(Equal(domain.User{}))
					Expect(response.Success).To(Equal(false))
				})
			})

			Context("when adding one invalid user (invalid email)", func() {

				var newUser *domain.NewUser
				var response users.CreateUserResponse_v0

				BeforeEach(func() {

					newUser = gory.Build("newUserInvalidEmail").(*domain.NewUser)
					sendRequestHelper("POST", "/api/users", users.CreateUserRequest_v0{
						User: *newUser,
					}, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns nil user", func() {
					Expect(response.User).To(Equal(domain.User{}))
					Expect(response.Success).To(Equal(false))
				})
			})
		})
		Context("when API user is anonymous", func() {

			var newUser *domain.NewUser
			var response users.CreateUserResponse_v0

			BeforeEach(func() {

				newUser = gory.Build("newUser").(*domain.NewUser)
				sendRequestHelper("POST", "/api/users", users.CreateUserRequest_v0{
					User: *newUser,
				}, nil, &response)

			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns newly-created user", func() {
				Expect(response.User.Email).To(Equal(newUser.Email))
				Expect(response.User.Roles).To(BeNil())
				Expect(response.User.Status).To(Equal(domain.StatusPending))
				Expect(response.Success).To(Equal(true))
			})
		})
		Context("when API user is an active user", func() {

			var activeUser *domain.User
			var newUser *domain.NewUser
			var response users.CreateUserResponse_v0

			BeforeEach(func() {

				activeUser = gory.Build("user").(*domain.User)
				db.Insert(repositories.UsersCollection, activeUser)

				newUser = gory.Build("newUser").(*domain.NewUser)
				sendRequestHelper("POST", "/api/users", users.CreateUserRequest_v0{
					User: *newUser,
				}, activeUser, &response)
			})

			It("returns status code of StatusForbidden (403)", func() {
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
		})

	})

	Describe("PUT /api/users", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})
			Context("when Action=`delete`", func() {

				var user1 *domain.User
				var user2 *domain.User

				BeforeEach(func() {
					// insert two valid users into database
					user1 = gory.Build("user").(*domain.User)
					user2 = gory.Build("user").(*domain.User)
					db.Insert(repositories.UsersCollection, user1)
					db.Insert(repositories.UsersCollection, user2)
				})

				Context("when delete one of the users", func() {
					var response users.UpdateUsersResponse_v0
					BeforeEach(func() {
						sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{
							Action: "delete",
							IDs:    []string{user1.ID.Hex()},
						}, adminUser, &response)
					})

					It("returns status code of StatusOK (200)", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})
					It("returns OK", func() {
						Expect(response.Action).To(Equal("delete"))
						Expect(response.IDs).To(Equal([]string{user1.ID.Hex()}))
						Expect(response.Success).To(Equal(true))
					})
				})

				Context("when IDs array is empty", func() {
					var response users.UpdateUsersResponse_v0
					BeforeEach(func() {
						sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{
							Action: "delete",
							IDs:    []string{},
						}, adminUser, &response)
					})
					It("returns status code of StatusOK (200)", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})
					It("returns OK", func() {
						Expect(response.Action).To(Equal("delete"))
						Expect(response.IDs).To(BeNil())
						Expect(response.Success).To(Equal(true))
					})
				})

				Context("when one of the IDs is not a valid ObjectId", func() {
					var response users.UpdateUsersResponse_v0
					BeforeEach(func() {
						sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{
							Action: "delete",
							IDs:    []string{"INVALIDID"},
						}, adminUser, &response)
					})
					It("returns status code of StatusOK (200)", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})
					It("returns OK", func() {
						Expect(response.Action).To(Equal("delete"))
						Expect(response.IDs).To(Equal([]string{"INVALIDID"}))
						Expect(response.Success).To(Equal(true))
					})
				})

			})

			Context("when Action is not supported`", func() {

				var response users.UpdateUsersResponse_v0

				BeforeEach(func() {

					sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{
						Action: "NOTSUPPORTED",
						IDs:    []string{},
					}, adminUser, &response)

				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns error", func() {
					Expect(response.Action).To(Equal("NOTSUPPORTED"))
					Expect(response.IDs).To(BeNil())
					Expect(response.Success).To(Equal(false))
				})

			})

			Context("when PUTing a malformed JSON body", func() {

				var response users.UpdateUsersResponse_v0

				BeforeEach(func() {
					sendRequestHelper("PUT", "/api/users", "BADJSON", adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns failed response", func() {
					Expect(response.Success).To(Equal(false))
				})
			})
		})

		Context("when API user is anonymous", func() {

			var response users.UpdateUsersResponse_v0

			BeforeEach(func() {
				sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{}, nil, &response)
			})

			It("returns status code of StatusForbidden (403)", func() {
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
		})
		Context("when API user is an active user", func() {

			var activeUser *domain.User
			var response users.CreateUserResponse_v0

			BeforeEach(func() {

				activeUser = gory.Build("user").(*domain.User)
				db.Insert(repositories.UsersCollection, activeUser)

				sendRequestHelper("PUT", "/api/users", users.UpdateUsersRequest_v0{}, activeUser, &response)
			})

			It("returns status code of StatusForbidden (403)", func() {
				Expect(recorder.Code).To(Equal(http.StatusForbidden))
			})
		})
	})

	Describe("DELETE /api/users", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})
			Context("when no users exists", func() {

				var response users.DeleteAllUsersResponse_v0

				BeforeEach(func() {

					sendRequestHelper("DELETE", "/api/users", nil, adminUser, &response)

				})

				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns OK", func() {
					Expect(response.Success).To(Equal(true))
				})

			})

			Context("when two users exist", func() {

				var user1 *domain.User
				var user2 *domain.User

				var response users.DeleteAllUsersResponse_v0

				BeforeEach(func() {

					// insert two valid users into database
					user1 = gory.Build("user").(*domain.User)
					user2 = gory.Build("user").(*domain.User)
					db.Insert(repositories.UsersCollection, user1)
					db.Insert(repositories.UsersCollection, user2)

					sendRequestHelper("DELETE", "/api/users", nil, adminUser, &response)

				})

				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns nil user", func() {
					Expect(response.Success).To(Equal(true))
				})

			})
		})

	})

	Describe("GET /api/users/{id}", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})
			Context("when user exists", func() {

				var user *domain.User
				var response users.GetUserResponse_v0
				BeforeEach(func() {

					// insert a user into database
					user = gory.Build("user").(*domain.User)
					db.Insert(repositories.UsersCollection, user)

					sendRequestHelper("GET", fmt.Sprintf("/api/users/%v", user.ID.Hex()), nil, adminUser, &response)
				})

				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns specified user", func() {
					Expect(response.User.ID).To(Equal(user.ID))
					Expect(response.User.Email).To(Equal(user.Email))
					Expect(response.Success).To(Equal(true))
				})
				It("should not return private fields", func() {
					Expect(response.User.ConfirmationCode).To(Equal(""))
					Expect(response.User.HashedPassword).To(Equal(""))
				})
			})

			Context("when user does not exists", func() {

				var response users.GetUserResponse_v0
				BeforeEach(func() {
					sendRequestHelper("GET", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), nil, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns unsuccessful response", func() {
					Expect(response.User).To(Equal(domain.User{}))
					Expect(response.Success).To(Equal(false))
				})
			})

			Context("when user `id` is invalid", func() {

				var response users.GetUserResponse_v0
				BeforeEach(func() {
					sendRequestHelper("GET", "/api/users/INVALIDID", nil, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns unsuccessful response", func() {
					Expect(response.User).To(Equal(domain.User{}))
					Expect(response.Success).To(Equal(false))
				})
			})
		})

	})

	Describe("GET /api/users/{id}/confirm", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})

			Context("when user exists", func() {

				var user *domain.User

				Context("when status is domain.StatusPending", func() {

					BeforeEach(func() {
						// insert a user into database
						user = gory.Build("userUnconfirmed").(*domain.User)
						user.GenerateConfirmationCode()
						db.Insert(repositories.UsersCollection, user)
					})

					Context("when code is correct", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode), nil, adminUser, &response)
						})
						It("returns status code of StatusOK (200)", func() {
							Expect(recorder.Code).To(Equal(http.StatusOK))
						})
						It("returns OK", func() {
							Expect(response.User.ID).To(Equal(user.ID))
							Expect(response.Success).To(Equal(true))
							Expect(response.Code).To(Equal(user.ConfirmationCode))
						})
					})
					Context("when code is incorrect", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=WRONGCODE", user.ID.Hex()), nil, adminUser, &response)
						})
						It("returns status code of StatusBadRequest (400)", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})
						It("returns not OK", func() {
							Expect(response.Success).To(Equal(false))
						})

					})
					Context("when code is empty/unspecified", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", user.ID.Hex()), nil, adminUser, &response)
						})
						It("returns status code of StatusBadRequest (400)", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})
						It("returns not OK", func() {
							Expect(response.Success).To(Equal(false))
						})

					})
				})

				Context("when status is not domain.StatusPending", func() {
					BeforeEach(func() {
						// insert a user into database
						user = gory.Build("user").(*domain.User)
						db.Insert(repositories.UsersCollection, user)

					})
					Context("when code is correct", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode), nil, adminUser, &response)
						})
						It("returns status code of StatusBadRequest (400)", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})
						It("returns not OK", func() {
							Expect(response.Success).To(Equal(false))
						})
					})
					Context("when code is incorrect", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=WRONGCODE", user.ID.Hex()), nil, adminUser, &response)
						})
						It("returns status code of StatusBadRequest (400)", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})
						It("returns not OK", func() {
							Expect(response.Success).To(Equal(false))
						})

					})
					Context("when code is empty/unspecified", func() {
						var response users.ConfirmUserResponse_v0
						BeforeEach(func() {
							sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", user.ID.Hex()), nil, adminUser, &response)
						})
						It("returns status code of StatusBadRequest (400)", func() {
							Expect(recorder.Code).To(Equal(http.StatusBadRequest))
						})
						It("returns not OK", func() {
							Expect(response.Success).To(Equal(false))
						})
					})
				})
			})

			Context("when user does not exist", func() {
				var response users.GetUserResponse_v0
				BeforeEach(func() {
					sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", bson.NewObjectId().Hex()), nil, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns unsuccessful response", func() {
					Expect(response.Success).To(Equal(false))
				})
			})

		})
	})

	Describe("PUT /api/users/{id}", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})

			Context("when user exists", func() {
				var user *domain.User
				var changedFields *domain.User
				var response users.UpdateUserResponse_v0

				Context("when all user fields are non-empty", func() {

					BeforeEach(func() {
						// insert a user into database
						user = gory.Build("user").(*domain.User)
						db.Insert(repositories.UsersCollection, user)

						loc, _ := time.LoadLocation("")
						changedFields = &domain.User{
							ID:               bson.NewObjectId(),
							Username:         "username",
							Email:            "new@email.com",
							Roles:            domain.Roles{domain.RoleAdmin},
							Status:           domain.StatusSuspended,
							LastModifiedDate: time.Date(1975, time.October, 1, 1, 1, 1, 1, loc),
							CreatedDate:      time.Date(1975, time.July, 1, 1, 1, 1, 1, loc),
						}
						// send request
						sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), users.UpdateUserRequest_v0{
							User: *changedFields,
						}, adminUser, &response)

					})

					It("returns status code of StatusOK (200)", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})
					It("returns successful response", func() {
						Expect(response.Success).To(Equal(true))
					})
					It("should only update permissible fields", func() {
						Expect(response.User.Username).To(Equal(changedFields.Username))
						Expect(response.User.Email).To(Equal(changedFields.Email))
						Expect(response.User.Roles).To(Equal(changedFields.Roles))
						Expect(response.User.Status).To(Equal(changedFields.Status))
					})

					It("should not update non-permissible fields", func() {
						Expect(response.User.ID).NotTo(Equal(changedFields.ID))
						Expect(response.User.LastModifiedDate).NotTo(Equal(changedFields.LastModifiedDate))
						Expect(response.User.CreatedDate).NotTo(Equal(changedFields.CreatedDate))
					})
				})
				Context("when all fields are empty", func() {

					BeforeEach(func() {
						// insert a user into database
						user = gory.Build("user").(*domain.User)
						db.Insert(repositories.UsersCollection, user)

						changedFields = &domain.User{}
						// send request
						sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), users.UpdateUserRequest_v0{
							User: *changedFields,
						}, adminUser, &response)

					})

					It("returns status code of StatusOK (200)", func() {
						Expect(recorder.Code).To(Equal(http.StatusOK))
					})
					It("returns successful response", func() {
						Expect(response.Success).To(Equal(true))
					})
					It("should not update fields", func() {

						// check against changedFields
						Expect(response.User.Username).NotTo(Equal(changedFields.Username))
						Expect(response.User.Email).NotTo(Equal(changedFields.Email))
						Expect(response.User.Roles).NotTo(Equal(changedFields.Roles))
						Expect(response.User.Status).NotTo(Equal(changedFields.Status))
						Expect(response.User.ID).NotTo(Equal(changedFields.ID))
						Expect(response.User.LastModifiedDate).NotTo(Equal(changedFields.LastModifiedDate))
						Expect(response.User.CreatedDate).NotTo(Equal(changedFields.CreatedDate))

						// check with original user object
						Expect(response.User.Username).To(Equal(user.Username))
						Expect(response.User.Email).To(Equal(user.Email))
						Expect(response.User.Roles).To(Equal(user.Roles))
						Expect(response.User.Status).To(Equal(user.Status))
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.User.LastModifiedDate).NotTo(Equal(user.LastModifiedDate))
						Expect(response.User.CreatedDate).NotTo(Equal(user.CreatedDate))
					})
				})

				Context("when PUTing a malformed JSON body", func() {

					BeforeEach(func() {
						user = gory.Build("user").(*domain.User)
						db.Insert(repositories.UsersCollection, user)

						sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), "BADJSON", adminUser, &response)
					})

					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns failed response", func() {
						Expect(response.Success).To(Equal(false))
					})
				})

			})

			Context("when user does not exist", func() {
				var changedFields *domain.User
				var response users.UpdateUserResponse_v0
				BeforeEach(func() {
					// insert a user into
					changedFields = &domain.User{}
					// send request
					sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), users.UpdateUserRequest_v0{
						User: *changedFields,
					}, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns successful response", func() {
					Expect(response.Success).To(Equal(false))
				})
			})

			Context("when user id is invalid", func() {
				var changedFields *domain.User
				var response users.UpdateUserResponse_v0
				BeforeEach(func() {
					// insert a user into
					changedFields = &domain.User{}
					// send request
					sendRequestHelper("PUT", "/api/users/NOTANID1a548b7539d00001f", users.UpdateUserRequest_v0{
						User: *changedFields,
					}, adminUser, &response)
				})

				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns successful response", func() {
					Expect(response.Success).To(Equal(false))
				})
			})

		})
	})

	Describe("DELETE /api/users/{id}", func() {

		Context("when API user is an active admin", func() {
			var adminUser *domain.User

			BeforeEach(func() {
				adminUser = gory.Build("adminAPIUser").(*domain.User)
				db.Insert(repositories.UsersCollection, adminUser)
			})
			Context("when user exists", func() {
				var user1 *domain.User
				var user2 *domain.User

				var response users.DeleteUserResponse_v0

				BeforeEach(func() {

					// insert two valid users into database
					user1 = gory.Build("user").(*domain.User)
					user2 = gory.Build("user").(*domain.User)
					db.Insert(repositories.UsersCollection, user1)
					db.Insert(repositories.UsersCollection, user2)

					sendRequestHelper("DELETE", fmt.Sprintf("/api/users/%v", user1.ID.Hex()), nil, adminUser, &response)

				})
				It("returns status code of StatusOK (200)", func() {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				})
				It("returns successful response", func() {
					Expect(response.Success).To(Equal(true))
				})
			})

			Context("when user does not exists", func() {

				var response users.DeleteUserResponse_v0

				BeforeEach(func() {
					sendRequestHelper("DELETE", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), nil, adminUser, &response)
				})
				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns successful response", func() {
					Expect(response.Success).To(Equal(false))
				})
			})

			Context("when user id is invalid", func() {

				var response users.DeleteUserResponse_v0

				BeforeEach(func() {
					sendRequestHelper("DELETE", "/api/users/INVALID", nil, adminUser, &response)
				})
				It("returns status code of StatusBadRequest (400)", func() {
					Expect(recorder.Code).To(Equal(http.StatusBadRequest))
				})
				It("returns successful response", func() {
					Expect(response.Success).To(Equal(false))
				})

			})
		})
	})
})
