package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares/context"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Middleware Tests", func() {
	Describe("ContextHandlerFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)
				ctx := context.New()

				var handler domain.ContextHandlerFunc
				handler = func(w http.ResponseWriter, req *http.Request, ctx domain.IContext) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
				}
				handler.ServeHTTP(recorder, request, ctx)

			})

		})

	})

	Describe("MiddlewareFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)

				next := func(w http.ResponseWriter, req *http.Request) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
				}

				var handler domain.MiddlewareFunc
				handler = func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
					next(w, req)
				}
				handler.ServeHTTP(recorder, request, next)

			})

		})

	})

	Describe("ContextMiddlewareFunc Type", func() {
		Describe("ServeHTTP()", func() {
			It("should be working", func() {
				recorder := httptest.NewRecorder()
				request, _ := http.NewRequest("GET", "/api/test", nil)
				ctx := context.New()

				next := func(w http.ResponseWriter, req *http.Request) {
					val := ctx.Get(req, "TESTKEY").(string)

					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))
					Expect(val).To(Equal("TESTVALUE"))
				}

				var handler domain.ContextMiddlewareFunc
				handler = func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc, ctx domain.IContext) {
					Expect(req.URL.Path).To(Equal("/api/test"))
					Expect(req.Method).To(Equal("GET"))

					ctx.Set(req, "TESTKEY", "TESTVALUE")

					next(w, req)
				}
				handler.ServeHTTP(recorder, request, next, ctx)

			})

		})

	})
})
