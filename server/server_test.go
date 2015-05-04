package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
)

var _ = Describe("Server", func() {
	var server *Server
	var routes *Routes

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Describe("Default server config", func() {

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
		server = NewServer(&Config{
			Database: &DatabaseOptions{
				ServerName:   TestDatabaseServerName,
				DatabaseName: TestDatabaseName,
			},
			Renderer: &RendererOptions{},
			Routes:   routes,
		})
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
