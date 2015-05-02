package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"github.com/unrolled/render"
	"gopkg.in/mgo.v2"
	"net/http"
)

var _ = Describe("Context", func() {

	var request *http.Request

	BeforeEach(func() {
		request, _ = http.NewRequest("GET", "/test", nil)
	})

	Describe("Database", func() {

		var session *mgo.Session
		var db *Database

		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/test", nil)
			session, _ = mgo.Dial("localhost")
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
			r = &Renderer{render.New(render.Options{})}
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
