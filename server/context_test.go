package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
)

var _ = Describe("Context", func() {

	var request *http.Request

	BeforeEach(func() {
		request, _ = http.NewRequest("GET", "/test", nil)
	})

	Describe("Database", func() {

		var session *DatabaseSession
		var db *Database

		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/test", nil)
			session = NewSession(DatabaseOptions{
				ServerName:   TestDatabaseServerName,
				DatabaseName: TestDatabaseName,
			})
			db = &Database{session.DB("test-db")}
		})

		Context("when db is a valid object", func() {
			It("returns original object", func() {
				SetDbCtx(request, db)

				retDb := DbCtx(request)
				Expect(retDb).To(Equal(db))
			})
		})

		Context("when db does not exist in context", func() {
			It("returns original object", func() {

				retDb := DbCtx(request)
				Expect(retDb).To(BeNil())
			})
		})

		AfterEach(func() {
			session.Close()
		})
	})

	Describe("Render", func() {

		var r *Renderer

		BeforeEach(func() {
			r = NewRenderer(RendererOptions{})
		})
		Context("when render is a valid object", func() {
			It("returns original object", func() {
				SetRendererCtx(request, r)

				renderer := RendererCtx(request)
				Expect(renderer).To(Equal(r))
			})
		})

		Context("when render does not exist in context", func() {
			It("returns original object", func() {

				renderer := RendererCtx(request)
				Expect(renderer).To(BeNil())
			})
		})
	})

	Describe("ClearHandler()", func() {

	})
})
