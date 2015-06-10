package libs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/libs"
	"net/http"
)

var _ = Describe("Merge Map Structs Tests", func() {

	Describe("MergeACLMap()", func() {
		stub := func(user *domain.User, req *http.Request, ctx domain.IContext) (bool, string) {
			return true, ""
		}
		firstMap := domain.ACLMap{
			"first": stub,
		}
		secondMap := domain.ACLMap{
			"second": stub,
		}
		var result domain.ACLMap
		BeforeEach(func() {
			result = libs.MergeACLMap(&firstMap, &secondMap)
		})
		It("should return a new map", func() {
			Expect(result["first"]).ToNot(BeNil())
		})
		It("should return a new map", func() {
			Expect(result["second"]).ToNot(BeNil())
		})
		It("should return a new map", func() {
			Expect(result["third"]).To(BeNil())
		})
	})
	Describe("MergeRoutes()", func() {
		firstRoutes := domain.Routes{
			domain.Route{},
		}
		secondRoutes := domain.Routes{
			domain.Route{},
		}
		thirdRoutes := domain.Routes{
			domain.Route{},
		}
		var result domain.Routes
		It("should return a new map", func() {
			result = libs.MergeRoutes(&firstRoutes, &secondRoutes)
			Expect(len(result)).To(Equal(2))
		})
		It("should return a new map", func() {
			result = libs.MergeRoutes(&firstRoutes, &secondRoutes)
			result = libs.MergeRoutes(&result, &thirdRoutes)
			Expect(len(result)).To(Equal(3))
		})
	})
})
