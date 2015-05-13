package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"github.com/sogko/golang-rest-api-server-example/middlewares"
	"github.com/sogko/golang-rest-api-server-example/server"
	"net/http"
)

var _ = Describe("Server", func() {
	var s *server.Server
	var routes *server.Routes

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Describe("Default server config", func() {

		routes = &server.Routes{
			server.Route{"Test", "GET", "/api/test", "0.1", server.RouteHandlers{
				"0.1": func(rw http.ResponseWriter, req *http.Request, ctx domain.IContext) {
					r := ctx.GetRendererCtx(req)
					r.JSON(rw, http.StatusOK, map[string]string{
						"ok": "ok",
					})
				},
			}},
		}
		s = server.NewServer(&server.Config{
			Database: &middlewares.MongoDBOptions{
				ServerName:   TestDatabaseServerName,
				DatabaseName: TestDatabaseName,
			},
			Renderer:       &middlewares.RendererOptions{},
			Routes:         routes,
			TokenAuthority: &middlewares.TokenAuthorityOptions{},
		})
	})

	Describe("Bad database server config", func() {
		It("should panic", func(done Done) {
			Expect(func() {

				db := middlewares.MongoDB{}
				db.NewSession(middlewares.MongoDBOptions{
					ServerName:   "BadServerName",
					DatabaseName: "BadDBName",
					DialTimeout:  1,
				})
			}).Should(Panic())
			close(done)
		}, 15)
	})
})
