package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"net/http"
)

var _ = Describe("Server", func() {
	var s *server.Server
	var routes *domain.Routes

	BeforeEach(func() {
	})

	AfterEach(func() {
	})

	Describe("Default server config", func() {

		routes = &domain.Routes{
			domain.Route{"Test", "GET", "/api/test", "0.1", domain.RouteHandlers{
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
		}).SetupRoutes()
	})

	Describe("Bad database server config", func() {
		It("should panic", func(done Done) {
			Expect(func() {

				db := middlewares.NewMongoDB(&middlewares.MongoDBOptions{
					ServerName:   "BadServerName",
					DatabaseName: "BadDBName",
					DialTimeout:  1,
				})
				db.NewSession()
			}).Should(Panic())
			close(done)
		}, 15)
	})
})
