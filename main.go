package main

import (
	"github.com/sogko/golang-rest-api-server-example/server"
)

func main() {

	// initialize server components
	components := server.Components{
		DatabaseSession: server.NewSession(server.DatabaseOptions{
			ServerName:   "localhost",
			DatabaseName: "test-app",
		}),
		Renderer: server.NewRenderer(server.RendererOptions{
			IndentJSON: true,
		}),
	}

	// init server and run
	s := server.NewServer(&components)
	s.Run(":3001")

}
