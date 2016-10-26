package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"net/http"
)

var _ = Describe("AccessController Tests", func() {
	Describe("ACLMap Struct", func() {
		Describe("Append()", func() {
			stub := func(req *http.Request, user domain.IUser) (bool, string) {
				return true, ""
			}
			firstMap := domain.ACLMap{
				"first": stub,
			}
			secondMap := domain.ACLMap{
				"second": stub,
			}
			thirdMap := domain.ACLMap{
				"third": stub,
			}
			var result domain.ACLMap
			var result2 domain.ACLMap
			BeforeEach(func() {
				result = firstMap.Append(&secondMap)
				result2 = firstMap.Append(&secondMap, &thirdMap)
			})
			It("should return a new map", func() {
				Expect(result["first"]).ToNot(BeNil())
				Expect(result2["first"]).ToNot(BeNil())
			})
			It("should return a new map", func() {
				Expect(result["second"]).ToNot(BeNil())
				Expect(result2["second"]).ToNot(BeNil())
			})
			It("should return a new map", func() {
				Expect(result["third"]).To(BeNil())
				Expect(result2["second"]).ToNot(BeNil())
			})
		})
	})
})
