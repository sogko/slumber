package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/utils"
)

var _ = Describe("MapFromJSON", func() {

	It("should map JSON string bytes to map[] if data is a valid JSON", func() {
		data := []byte(`
		{
			"a": "isString",
			"b": 100,
			"c": true
		}
		`)
		body := MapFromJSON(data)
		Expect(body["a"]).To(Equal("isString"))
		Expect(body["b"]).To(Equal(float64(100)))
		Expect(body["c"]).To(Equal(true))

	})

	It("should panic if data is an invalid json ", func() {
		data := []byte("{this is an invalid json}")
		Expect(func() {
			_ = MapFromJSON(data)
		}).Should(Panic())
	})
})
