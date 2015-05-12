package main

import (
	"github.com/sogko/golang-rest-api-server-example/middlewares"
	"github.com/sogko/golang-rest-api-server-example/server"
	"io/ioutil"
	"errors"
	"fmt"
)

func main() {

	// try to load signing keys for token authority
	privateSigningKey, err := ioutil.ReadFile("keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err := ioutil.ReadFile("keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	// load routes
	routes := GetRoutes()

	// set server configuration
	config := server.Config{
		Database: &middlewares.MongoDBOptions{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		},
		Renderer: &middlewares.RendererOptions{
			IndentJSON: true,
		},
		Routes: routes,
		TokenAuthority: &middlewares.TokenAuthorityOptions{
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		},
	}

	// init server and run
	s := server.NewServer(&config)
	// bam!
	s.Run(":3001")
}