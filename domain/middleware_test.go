package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Middleware Tests", func() {
	Describe("ContextHandlerFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {

				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)

				var handler domain.ContextHandlerFunc
				handler = func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
				}
				handler.ServeHTTP(recorder, request, nil)

			})

		})

	})
})
