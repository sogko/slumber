package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
)

var _ = Describe("Routes Tests", func() {
	Describe("Routes Struct", func() {
		Describe("Append()", func() {
			firstRoutes := domain.Routes{
				domain.Route{},
			}
			secondRoutes := domain.Routes{
				domain.Route{},
			}
			thirdRoutes := domain.Routes{
				domain.Route{},
			}
			It("should return a new map", func() {
				result := firstRoutes.Append(&secondRoutes)
				Expect(len(result)).To(Equal(2))
			})
			It("should return a new map", func() {
				result := firstRoutes.Append(&secondRoutes)
				result = result.Append(&thirdRoutes)
				Expect(len(result)).To(Equal(3))
			})
			It("should return a new map", func() {
				result := firstRoutes.Append(&secondRoutes, &thirdRoutes)
				Expect(len(result)).To(Equal(3))
			})
		})
	})
})
