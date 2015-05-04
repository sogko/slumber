package main

import (
	"github.com/sogko/golang-rest-api-server-example/server"
)

func main() {

	// load routes
	routes := GetRoutes()

	// set server configuration
	config := server.Config{
		Database: &server.DatabaseOptions{
			ServerName:   "localhost",
			DatabaseName: "test-go-app",
		},
		Renderer: &server.RendererOptions{
			IndentJSON: true,
		},
		Routes: routes,
	}

	// init server and run
	s := server.NewServer(&config)
	// bam!
	s.Run(":3001")
}
