package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
)

var _ = Describe("Server", func() {
	var server *Server
	var session *DatabaseSession
	var routes *Routes
	var renderer *Renderer

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Describe("Default server config", func() {

		// set up server with test components
		session = NewSession(DatabaseOptions{
			ServerName:   TestDatabaseServerName,
			DatabaseName: TestDatabaseName,
		})

		renderer = NewRenderer(RendererOptions{})

		routes = &Routes{
			Route{"Test", "GET", "/api/test", "0.1", RouteHandlers{
				"0.1": func(rw http.ResponseWriter, req *http.Request) {
					r := RendererCtx(req)
					r.JSON(rw, http.StatusOK, map[string]string{
						"ok": "ok",
					})
				},
			}},
		}
		components := Components{
			DatabaseSession: session,
			Renderer:        renderer,
			Routes:          routes,
		}
		server = NewServer(&components)
	})

	Describe("Bad database server config", func() {
		It("should panic", func(done Done) {
			Expect(func() {
				NewSession(DatabaseOptions{
					ServerName:   "BadServerName",
					DatabaseName: "BadDBName",
					DialTimeout:  1,
				})
			}).Should(Panic())
			close(done)
		}, 15)
	})
})
