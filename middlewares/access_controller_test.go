package middlewares_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("AccessController", func() {

	var request *http.Request
	var ctx *middlewares.Context
	var adminUser *domain.User
	var normalUser *domain.User
	var anotherNormalUser *domain.User
	var ac *middlewares.AccessController
	var aclMap domain.ACLMap

	BeforeEach(func() {

		// dummy request and context object
		request, _ = http.NewRequest("GET", "/test/api", nil)
		ctx = middlewares.NewContext()

		// create users with roles
		adminUser = &domain.User{
			ID:    bson.NewObjectId(),
			Roles: domain.Roles{domain.RoleAdmin},
		}
		normalUser = &domain.User{
			ID:    bson.NewObjectId(),
			Roles: domain.Roles{domain.RoleUser},
		}
		anotherNormalUser = &domain.User{
			ID:    bson.NewObjectId(),
			Roles: domain.Roles{domain.RoleUser},
		}

		// create test ACL Map
		aclMap = domain.ACLMap{
			"ListUsers": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
				// does not require special privileges
				return true, ""
			},
			"EditUser": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
				// This is an example of determining access to current resource by storing/retrieving
				// contextual data from IContext
				// In practical use, the handler might use route params (for eg /api/users/{id})
				// to get current resource context
				currObj := ctx.GetCurrentObjectCtx(req)
				if currObj == nil {
					currObj = &domain.User{}
				}
				userObj := currObj.(*domain.User)
				return user.HasRole(domain.RoleAdmin) || user.ID == userObj.ID, ""
			},
		}

		// create and init AccessController
		ac = middlewares.NewAccessController()

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
				stub := func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
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
				ctx.SetCurrentObjectCtx(request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", adminUser)
				Expect(result).To(BeTrue())
			})
		})
		Context("when user is authorized (owns `user` resource)", func() {
			It("return OK", func() {
				ctx.SetCurrentObjectCtx(request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", normalUser)
				Expect(result).To(BeTrue())
			})
		})
		Context("when user is not authorized", func() {
			It("return not OK", func() {
				ctx.SetCurrentObjectCtx(request, anotherNormalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "EditUser", normalUser)
				Expect(result).To(BeFalse())
			})
		})
		Context("when action does not exists", func() {
			It("return not OK", func() {
				ctx.SetCurrentObjectCtx(request, normalUser)
				result, _ := ac.IsHTTPRequestAuthorized(request, ctx, "NonExistent", normalUser)
				Expect(result).To(BeFalse())
			})
		})
	})

	type TestResponse struct {
		Value string `json:"value,omitempty"`
		Success bool `json:"success,omitempty"`
		Message string  `json:"message,omitempty"`
	}
	Describe("Handler()", func() {
		// setup other upstream middlewares
		ctx := middlewares.NewContext()
		renderer := middlewares.NewRenderer(&middlewares.RendererOptions{
			IndentJSON: true,
		})

		Context("when request is authorized", func () {
			It("should be working", func() {

				// add test ACL map
				ac.Add(&domain.ACLMap{
					"TestAuthorized": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
						return true, ""
					},
				})

				recorder := httptest.NewRecorder()

				acHandler := ctx.Inject(ac.Handler("TestAuthorized", func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {}))
				renderer.Handler(recorder, request, acHandler, ctx)

				acHandler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusOK))

			})
		})

		Context("when request is forbidden", func () {
			It("should be working", func() {
				// add test ACL map
				ac.Add(&domain.ACLMap{
					"TestForbidden": func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
						return false, ""
					},
				})

				recorder := httptest.NewRecorder()

				acHandler := ctx.Inject(ac.Handler("TestForbidden", func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {}))
				renderer.Handler(recorder, request, acHandler, ctx)

				acHandler.ServeHTTP(recorder, request)
				Expect(recorder.Code).To(Equal(http.StatusForbidden))

			})
		})


	})
})
