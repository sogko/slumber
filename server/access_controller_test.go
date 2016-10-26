package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber-users"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/middlewares/renderer"
	"github.com/sogko/slumber/server"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
)

func SetCurrentObjectCtx(ctx domain.IContext, req *http.Request, user *users.User) {
	ctx.Set(req, "TESTCURRENTOBJECT", user)
}
func GetCurrentObjectCtx(ctx domain.IContext, req *http.Request) *users.User {
	return ctx.Get(req, "TESTCURRENTOBJECT").(*users.User)
}

var _ = Describe("AccessController", func() {

	var request *http.Request
	var ctx domain.IContext
	var adminUser *users.User
	var normalUser *users.User
	var anotherNormalUser *users.User
	var ac *server.AccessController
	var aclMap domain.ACLMap

	BeforeEach(func() {

		// dummy request and context object
		request, _ = http.NewRequest("GET", "/test/api", nil)
		ctx = context.New()

		// create users with roles
		adminUser = &users.User{
			ID:    bson.NewObjectId(),
			Roles: users.Roles{users.RoleAdmin},
		}
		normalUser = &users.User{
			ID:    bson.NewObjectId(),
			Roles: users.Roles{users.RoleUser},
		}
		anotherNormalUser = &users.User{
			ID:    bson.NewObjectId(),
			Roles: users.Roles{users.RoleUser},
		}

		// create test ACL Map
		aclMap = domain.ACLMap{
			"ListUsers": func(req *http.Request, user domain.IUser) (bool, string) {
				// does not require special privileges
				return true, ""
			},
			"EditUser": func(req *http.Request, user domain.IUser) (bool, string) {
				// This is an example of determining access to current resource by storing/retrieving
				// contextual data from IContext
				// In practical use, the handler might use route params (for eg /api/users/{id})
				// to get current resource context
				currentObj := GetCurrentObjectCtx(ctx, req)
				return (user.HasRole(users.RoleAdmin) || user.GetID() == currentObj.ID.Hex()), ""
			},
		}

		renderer := renderer.New(&renderer.Options{}, renderer.JSON)

		// create and init AccessController
		ac = server.NewAccessController(ctx, renderer)

	})

	Describe("Add()", func() {
		Context("when nothing has been added yet", func() {
			BeforeEach(func() {
			})
			It("should be empty", func() {
				Expect(ac.ACLMap).To(BeEmpty())
			})
			It("should be not have HandlerFuncs", func() {
				Expect(ac.ACLMap["ListUsers"]).To(BeNil())
			})
		})
		Context("when nothing has been added yet", func() {
			BeforeEach(func() {
				ac.Add(&aclMap)
			})
			It("should not be empty", func() {
				Expect(ac.ACLMap).ToNot(BeEmpty())
			})
			It("should be have HandlerFuncs", func() {
				Expect(ac.ACLMap["ListUsers"]).ToNot(BeNil())
				Expect(ac.ACLMap["EditUser"]).ToNot(BeNil())
			})
		})
		Context("when it already have something", func() {
			BeforeEach(func() {
				stub := func(req *http.Request, user domain.IUser) (bool, string) {
					return true, ""
				}
				ac.Add(&domain.ACLMap{
					"ListAdmins": stub,
				})
				ac.Add(&aclMap)
				ac.Add(&domain.ACLMap{
					"UpdateAdmins": stub,
				})

			})
			It("should not be empty", func() {
				Expect(ac.ACLMap).ToNot(BeEmpty())
			})
			It("should be have HandlerFuncs", func() {
				Expect(ac.ACLMap["ListAdmins"]).ToNot(BeNil())
				Expect(ac.ACLMap["ListUsers"]).ToNot(BeNil())
				Expect(ac.ACLMap["EditUser"]).ToNot(BeNil())
				Expect(ac.ACLMap["UpdateAdmins"]).ToNot(BeNil())
			})
		})
	})

	Describe("HasAction()", func() {
		BeforeEach(func() {
			ac.Add(&aclMap)
		})
		It("should return true if action exists", func() {
			Expect(ac.HasAction("ListUsers")).To(BeTrue())
			Expect(ac.HasAction("EditUser")).To(BeTrue())
		})
		It("should return false if action does not exist", func() {
			Expect(ac.HasAction("NonExistent")).To(BeFalse())
		})
	})

	Describe("IsHTTPRequestAuthorized()", func() {
		BeforeEach(func() {
			ac.Add(&aclMap)
		})
		Context("when user is authorized (an admin)", func() {
			It("return OK", func() {
				SetCurrentObjectCtx(ctx, request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", adminUser)
				Expect(result).To(BeTrue())
			})
		})
		Context("when user is authorized (owns `user` resource)", func() {
			It("return OK", func() {
				SetCurrentObjectCtx(ctx, request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", normalUser)
				Expect(result).To(BeTrue())
			})
		})
		Context("when user is not authorized", func() {
			It("return not OK", func() {
				SetCurrentObjectCtx(ctx, request, anotherNormalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", normalUser)
				Expect(result).To(BeFalse())
			})
		})
		Context("when action does not exists", func() {
			It("return not OK", func() {
				SetCurrentObjectCtx(ctx, request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "NonExistent", normalUser)
				Expect(result).To(BeFalse())
			})
		})
	})

	type TestResponse struct {
		Value   string `json:"value,omitempty"`
		Success bool   `json:"success,omitempty"`
		Message string `json:"message,omitempty"`
	}
	Describe("Handler()", func() {
		// setup other upstream middlewares
		ctx := context.New()
		renderer := renderer.New(&renderer.Options{
			IndentJSON: true,
		}, renderer.JSON)

		Context("when request is authorized", func() {
			It("should be working", func() {

				// add test ACL map
				ac.Add(&domain.ACLMap{
					"TestAuthorized": func(req *http.Request, user domain.IUser) (bool, string) {
						return true, ""
					},
				})

				recorder := httptest.NewRecorder()

				acHandler := ac.NewContextHandler("TestAuthorized", func(w http.ResponseWriter, req *http.Request) {})
				renderer.Handler(recorder, request, acHandler, ctx)

				acHandler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))

			})
		})

		Context("when request is forbidden", func() {
			It("should be working", func() {
				// add test ACL map
				ac.Add(&domain.ACLMap{
					"TestForbidden": func(req *http.Request, user domain.IUser) (bool, string) {
						return false, ""
					},
				})

				recorder := httptest.NewRecorder()

				acHandler := ac.NewContextHandler("TestForbidden", func(w http.ResponseWriter, req *http.Request) {})
				renderer.Handler(recorder, request, acHandler, ctx)

				acHandler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusForbidden))

			})
		})

	})
})
