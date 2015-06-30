package test_helpers_test

import (
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/test_helpers"
)

var _ = Describe("Test Helpers", func() {

	Describe("NewTestResource", func() {

		It("should return nil routes when options.NilRoutes=true", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: true,
			})
			Expect(testResource.Routes()).To(BeNil())
		})
		It("should return routes when options.NilRoutes=false", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: false,
			})
			Expect(testResource.Routes()).ToNot(BeNil())
		})
		It("should return context", func() {
			ctx := context.New()
			testResource := test_helpers.NewTestResource(ctx, nil, &test_helpers.TestResourceOptions{
				NilRoutes: false,
			})
			Expect(testResource.Context()).To(Equal(ctx))
		})
	})
	Describe("MapFromJSON", func() {

		It("should map JSON string bytes to map[] if data is a valid JSON", func() {
			data := []byte(`
				{
					"a": "isString",
					"b": 100,
					"c": true
				}
			`)
			body := test_helpers.MapFromJSON(data)
			Expect(body["a"]).To(Equal("isString"))
			Expect(body["b"]).To(Equal(float64(100)))
			Expect(body["c"]).To(Equal(true))

		})

		It("should panic if data is an invalid json ", func() {
			data := []byte("{this is an invalid json}")
			Expect(func() {
				_ = test_helpers.MapFromJSON(data)
			}).Should(Panic())
		})
	})

	Describe("DecodeResponseToType", func() {

		type TestResponseType struct {
			A string `json:"a"`
			B int    `json:"b"`
			C bool   `json:"c"`
		}

		It("should map ResponseRecorder body data to target type if data is a valid JSON", func() {
			data := []byte(`
				{
					"a": "isString",
					"b": 100,
					"c": true
				}
			`)

			var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
			recorder.Body.Write(data)

			var responseType TestResponseType
			test_helpers.DecodeResponseToType(recorder, &responseType)

			Expect(responseType).To(Equal(TestResponseType{
				A: "isString",
				B: 100,
				C: true,
			}))
		})

		It("should not panic if data is an invalid json ", func() {
			data := []byte("{this is an invalid json}")

			var recorder *httptest.ResponseRecorder = httptest.NewRecorder()
			recorder.Body.Write(data)

			Expect(func() {
				var responseType TestResponseType
				test_helpers.DecodeResponseToType(recorder, &responseType)
			}).ShouldNot(Panic())
		})
	})
})
