package server_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber-sessions"
	"github.com/sogko/slumber-users"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/middlewares/mongodb"
	"github.com/sogko/slumber/middlewares/renderer"
	"github.com/sogko/slumber/server"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

var _ = Describe("Server", func() {
	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR TEST ONLY
	privateSigningKey, err := ioutil.ReadFile("../keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("../keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	Describe("Basic sanity test", func() {
		ctx := context.New()

		db := mongodb.New(&mongodb.Options{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		})
		_ = db.NewSession()

		// init renderer
		renderer := renderer.New(&renderer.Options{
			IndentJSON: true,
		}, renderer.JSON)

		// set up users resource
		usersResource := users.NewResource(ctx, &users.Options{
			Database: db,
			Renderer: renderer,
		})

		// set up sessions resource
		sessionsResource := sessions.NewResource(ctx, &sessions.Options{
			Database:              db,
			Renderer:              renderer,
			PrivateSigningKey:     privateSigningKey,
			PublicSigningKey:      publicSigningKey,
			UserRepositoryFactory: usersResource.UserRepositoryFactory,
		})

		// init server
		s := server.NewServer(&server.Config{
			Context: ctx,
		})

		// set up router
		ac := server.NewAccessController(ctx, renderer)
		router := server.NewRouter(s.Context, ac)

		// add REST resources to router
		router.AddResources(sessionsResource, usersResource)

		// add middlewares
		s.UseContextMiddleware(renderer)
		s.UseMiddleware(sessionsResource.NewAuthenticator())

		s.UseContextMiddleware(renderer)

		s.UseRouter(router)

		It("should serve request", func() {
			// run server and it shouldn't panic
			go s.Run(":8001", server.Options{
				Timeout: 1*time.Millisecond,
			})
			time.Sleep(100 * time.Millisecond)

			// serve some urls
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", "/api/sessions", nil)

			s.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusForbidden))

			// gracefully stops server
			s.Stop()

		})
	})
})
