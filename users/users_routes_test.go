package users_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example"
	. "github.com/sogko/golang-rest-api-server-example/server"
	. "github.com/sogko/golang-rest-api-server-example/users"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"time"
)

const RequestAcceptHeader = "application/json;version=0.0"

var _ = Describe("Users API - /api/users; version=0.0", func() {

	var server *Server
	var session *DatabaseSession
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	// sendRequestHelper helps
	// - creates an HTTP request
	// - set the right API version
	// - serves request
	// - decode the response into the desired interfac
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

	Describe("GET /api/users", func() {

		Context("when no users exist", func() {

			var response ListResponse_v0

			BeforeEach(func() {
				sendRequestHelper("GET", "/api/users", nil, &response)
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns zero users", func() {
				Expect(len(response.Users)).To(Equal(0))
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when two users exist", func() {

			var response ListResponse_v0

			BeforeEach(func() {
				// insert two valid users into database
				collection := session.DB(TestDatabaseName).C(UsersCollection)
				collection.Insert(gory.Build("user"))
				collection.Insert(gory.Build("user"))

				sendRequestHelper("GET", "/api/users", nil, &response)
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns two users", func() {
				Expect(len(response.Users)).To(Equal(2))
				Expect(response.Success).To(Equal(true))
			})

		})
	})

	Describe("POST /api/users", func() {

		Context("when adding one valid users", func() {

			var newUser *User
			var response CreateResponse_v0

			BeforeEach(func() {

				newUser = gory.Build("user").(*User)

				sendRequestHelper("POST", "/api/users", CreateRequest_v0{
					User: *newUser,
				}, &response)

			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns newly-created user", func() {
				Expect(response.User.Email).To(Equal(newUser.Email))
				Expect(response.User.Roles).To(BeNil())
				Expect(response.User.Status).To(Equal(StatusPending))
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when trying to specify Roles and Status", func() {

			var newUser *User
			var response CreateResponse_v0

			BeforeEach(func() {

				newUser = gory.Build("user").(*User)
				newUser.Status = StatusActive
				newUser.Roles = Roles{RoleAdmin}
				sendRequestHelper("POST", "/api/users", CreateRequest_v0{
					User: *newUser,
				}, &response)

			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns newly-created user", func() {
				Expect(response.User.Email).To(Equal(newUser.Email))
				Expect(response.User.Roles).To(BeNil())
				Expect(response.User.Status).To(Equal(StatusPending))
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when POSTing a malformed JSON body", func() {

			var response CreateResponse_v0

			BeforeEach(func() {
				sendRequestHelper("POST", "/api/users", "BADJSON", &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns nil user", func() {
				Expect(response.User).To(Equal(User{}))
				Expect(response.Success).To(Equal(false))
			})
		})

		Context("when adding one invalid user (invalid email)", func() {

			var newUser *User
			var response CreateResponse_v0

			BeforeEach(func() {

				newUser = gory.Build("userInvalidEmail").(*User)
				sendRequestHelper("POST", "/api/users", CreateRequest_v0{
					User: *newUser,
				}, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns nil user", func() {
				Expect(response.User).To(Equal(User{}))
				Expect(response.Success).To(Equal(false))
			})
		})

	})

	Describe("PUT /api/users", func() {

		Context("when Action=`delete`", func() {

			var user1 *User
			var user2 *User

			BeforeEach(func() {

				// insert two valid users into database
				collection := session.DB(TestDatabaseName).C(UsersCollection)
				user1 = gory.Build("user").(*User)
				user2 = gory.Build("user").(*User)
				collection.Insert(user1)
				collection.Insert(user2)
			})

			Context("when delete one of the users", func() {
				var response UpdateListResponse_v0
				BeforeEach(func() {
					sendRequestHelper("PUT", "/api/users", UpdateListRequest_v0{
						Action: "delete",
						IDs:    []string{user1.ID.Hex()},
					}, &response)
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
				var response UpdateListResponse_v0
				BeforeEach(func() {
					sendRequestHelper("PUT", "/api/users", UpdateListRequest_v0{
						Action: "delete",
						IDs:    []string{},
					}, &response)
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
				var response UpdateListResponse_v0
				BeforeEach(func() {
					sendRequestHelper("PUT", "/api/users", UpdateListRequest_v0{
						Action: "delete",
						IDs:    []string{"INVALIDID"},
					}, &response)
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

			var response UpdateListResponse_v0

			BeforeEach(func() {

				sendRequestHelper("PUT", "/api/users", UpdateListRequest_v0{
					Action: "NOTSUPPORTED",
					IDs:    []string{},
				}, &response)

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

			var response UpdateListResponse_v0

			BeforeEach(func() {
				sendRequestHelper("PUT", "/api/users", "BADJSON", &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns failed response", func() {
				Expect(response.Success).To(Equal(false))
			})
		})
	})

	Describe("DELETE /api/users", func() {

		Context("when no users exists", func() {

			var response DeleteAllResponse_v0

			BeforeEach(func() {

				sendRequestHelper("DELETE", "/api/users", nil, &response)

			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns OK", func() {
				Expect(response.Success).To(Equal(true))
			})

		})

		Context("when two users exist", func() {

			var user1 *User
			var user2 *User

			var response DeleteAllResponse_v0

			BeforeEach(func() {

				// insert two valid users into database
				collection := session.DB(TestDatabaseName).C(UsersCollection)
				user1 = gory.Build("user").(*User)
				user2 = gory.Build("user").(*User)
				collection.Insert(user1)
				collection.Insert(user2)

				sendRequestHelper("DELETE", "/api/users", nil, &response)

			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns nil user", func() {
				Expect(response.Success).To(Equal(true))
			})

		})

	})

	Describe("GET /api/users/{id}", func() {

		Context("when user exists", func() {

			var user *User
			var response GetResponse_v0
			BeforeEach(func() {

				// insert a user into database
				user = gory.Build("user").(*User)
				collection := session.DB(TestDatabaseName).C(UsersCollection)
				collection.Insert(user)

				sendRequestHelper("GET", fmt.Sprintf("/api/users/%v", user.ID.Hex()), nil, &response)
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

			var response GetResponse_v0
			BeforeEach(func() {
				sendRequestHelper("GET", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), nil, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns unsuccessful response", func() {
				Expect(response.User).To(Equal(User{}))
				Expect(response.Success).To(Equal(false))
			})
		})

		Context("when user `id` is invalid", func() {

			var response GetResponse_v0
			BeforeEach(func() {
				sendRequestHelper("GET", "/api/users/INVALIDID", nil, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns unsuccessful response", func() {
				Expect(response.User).To(Equal(User{}))
				Expect(response.Success).To(Equal(false))
			})
		})

	})

	Describe("GET /api/users/{id}/confirm", func() {

		Context("when user exists", func() {

			var user *User

			Context("when status is StatusPending", func() {

				BeforeEach(func() {
					// insert a user into database
					user = gory.Build("userUnconfirmed").(*User)
					user.GenerateConfirmationCode()
					collection := session.DB(TestDatabaseName).C(UsersCollection)
					collection.Insert(user)
				})

				Context("when code is correct", func() {
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode), nil, &response)
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
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=WRONGCODE", user.ID.Hex()), nil, &response)
					})
					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns not OK", func() {
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.Success).To(Equal(false))
						Expect(response.Code).To(Equal("WRONGCODE"))
					})

				})
				Context("when code is empty/unspecified", func() {
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", user.ID.Hex()), nil, &response)
					})
					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns not OK", func() {
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.Success).To(Equal(false))
						Expect(response.Code).To(Equal(""))
					})

				})
			})

			Context("when status is not StatusPending", func() {
				BeforeEach(func() {
					// insert a user into database
					user = gory.Build("user").(*User)
					collection := session.DB(TestDatabaseName).C(UsersCollection)
					collection.Insert(user)

				})
				Context("when code is correct", func() {
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=%v", user.ID.Hex(), user.ConfirmationCode), nil, &response)
					})
					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns not OK", func() {
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.Success).To(Equal(false))
						Expect(response.Code).To(Equal(user.ConfirmationCode))
					})
				})
				Context("when code is incorrect", func() {
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm?code=WRONGCODE", user.ID.Hex()), nil, &response)
					})
					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns not OK", func() {
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.Success).To(Equal(false))
						Expect(response.Code).To(Equal("WRONGCODE"))
					})

				})
				Context("when code is empty/unspecified", func() {
					var response ConfirmUserResponse_v0
					BeforeEach(func() {
						sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", user.ID.Hex()), nil, &response)
					})
					It("returns status code of StatusBadRequest (400)", func() {
						Expect(recorder.Code).To(Equal(http.StatusBadRequest))
					})
					It("returns not OK", func() {
						Expect(response.User.ID).To(Equal(user.ID))
						Expect(response.Success).To(Equal(false))
						Expect(response.Code).To(Equal(""))
					})
				})
			})
		})

		Context("when user does not exist", func() {
			var response GetResponse_v0
			BeforeEach(func() {
				sendRequestHelper("GET", fmt.Sprintf("/api/users/%v/confirm", bson.NewObjectId().Hex()), nil, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns unsuccessful response", func() {
				Expect(response.User).To(Equal(User{}))
				Expect(response.Success).To(Equal(false))
			})
		})

	})

	Describe("PUT /api/users/{id}", func() {
		Context("when user exists", func() {
			var user *User
			var changedFields *User
			var response UpdateResponse_v0

			Context("when all user fields are non-empty", func() {

				BeforeEach(func() {
					// insert a user into database
					user = gory.Build("user").(*User)
					collection := session.DB(TestDatabaseName).C(UsersCollection)
					collection.Insert(user)

					loc, _ := time.LoadLocation("")
					changedFields = &User{
						ID:               bson.NewObjectId(),
						Username:         "username",
						Email:            "new@email.com",
						Roles:            Roles{RoleAdmin},
						Status:           StatusSuspended,
						LastModifiedDate: time.Date(1975, time.October, 1, 1, 1, 1, 1, loc),
						CreatedDate:      time.Date(1975, time.July, 1, 1, 1, 1, 1, loc),
					}
					// send request
					sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), UpdateRequest_v0{
						User: *changedFields,
					}, &response)

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
					user = gory.Build("user").(*User)
					collection := session.DB(TestDatabaseName).C(UsersCollection)
					collection.Insert(user)

					changedFields = &User{}
					// send request
					sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), UpdateRequest_v0{
						User: *changedFields,
					}, &response)

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

				var response CreateResponse_v0

				BeforeEach(func() {
					sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", user.ID.Hex()), "BADJSON", &response)
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
			var changedFields *User
			var response UpdateResponse_v0
			BeforeEach(func() {
				// insert a user into
				changedFields = &User{}
				// send request
				sendRequestHelper("PUT", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), UpdateRequest_v0{
					User: *changedFields,
				}, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns successful response", func() {
				Expect(response.Success).To(Equal(false))
			})
		})

		Context("when user id is invalid", func() {
			var changedFields *User
			var response UpdateResponse_v0
			BeforeEach(func() {
				// insert a user into
				changedFields = &User{}
				// send request
				sendRequestHelper("PUT", "/api/users/INVALIDID", UpdateRequest_v0{
					User: *changedFields,
				}, &response)
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns successful response", func() {
				Expect(response.Success).To(Equal(false))
			})
		})

	})

	Describe("DELETE /api/users/{id}", func() {

		Context("when user exists", func() {
			var user1 *User
			var user2 *User

			var response DeleteResponse_v0

			BeforeEach(func() {

				// insert two valid users into database
				collection := session.DB(TestDatabaseName).C(UsersCollection)
				user1 = gory.Build("user").(*User)
				user2 = gory.Build("user").(*User)
				collection.Insert(user1)
				collection.Insert(user2)

				sendRequestHelper("DELETE", fmt.Sprintf("/api/users/%v", user1.ID.Hex()), nil, &response)

			})
			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns successful response", func() {
				Expect(response.Success).To(Equal(true))
			})
		})

		Context("when user does not exists", func() {

			var response DeleteResponse_v0

			BeforeEach(func() {
				sendRequestHelper("DELETE", fmt.Sprintf("/api/users/%v", bson.NewObjectId().Hex()), nil, &response)
			})
			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns successful response", func() {
				Expect(response.Success).To(Equal(false))
			})
		})

		Context("when user id is invalid", func() {

			var response DeleteResponse_v0

			BeforeEach(func() {
				sendRequestHelper("DELETE", "/api/users/INVALID", nil, &response)
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
