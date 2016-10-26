package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sogko/slumber-sessions"
	"github.com/sogko/slumber-users"
	"github.com/sogko/slumber/middlewares/context"
	"github.com/sogko/slumber/middlewares/mongodb"
	"github.com/sogko/slumber/middlewares/renderer"
	"github.com/sogko/slumber/server"
)

func main() {

	// try to load signing keys for token authority
	// NOTE: DO NOT USE THESE KEYS FOR PRODUCTION! FOR DEMO ONLY
	privateSigningKey, err := ioutil.ReadFile("keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	// Casting keys loaded to proper type
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateSigningKey)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error parsing private key: %v", err.Error())))
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicSigningKey)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error parsing public key: %v", err.Error())))
	}

	// create current project context
	ctx := context.New()

	// set up DB session
	db := mongodb.New(&mongodb.Options{
		ServerName:   "localhost",
		DatabaseName: "test-go-app",
	})
	_ = db.NewSession()

	// set up Renderer (unrolled_render)
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
		PrivateSigningKey:     privateKey,
		PublicSigningKey:      publicKey,
		Database:              db,
		Renderer:              renderer,
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
	s.UseMiddleware(sessionsResource.NewAuthenticator())

	// setup router
	s.UseRouter(router)

	// bam!
	s.Run(":3001", server.Options{
		Timeout: 10 * time.Second,
	})
}
