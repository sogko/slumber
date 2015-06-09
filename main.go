package main

import (
	"errors"
	"fmt"
	"github.com/sogko/slumber/libs"
	"github.com/sogko/slumber/middlewares"
	"github.com/sogko/slumber/server"
	"github.com/sogko/slumber/sessions"
	"github.com/sogko/slumber/users"
	"io/ioutil"
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

	// load routes
	routes := users.UsersAPIRoutes
	routes = libs.MergeRoutes(&routes, &sessions.SessionsAPIRoutes)

	// load ACL map
	aclMap := users.UsersAPIACL
	aclMap = libs.MergeACLMap(&aclMap, &sessions.SessionsAPIACL)

	// set server configuration
	config := server.Config{
		Database: &middlewares.MongoDBOptions{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		},
		Renderer: &middlewares.RendererOptions{
			IndentJSON: true,
		},
		Routes: &routes,
		TokenAuthority: &middlewares.TokenAuthorityOptions{
			PrivateSigningKey: privateSigningKey,
			PublicSigningKey:  publicSigningKey,
		},
		ACLMap: &aclMap,
	}

	// init server and run
	s := server.NewServer(&config).SetupRoutes()

	// bam!
	s.Run(":3001")
}
