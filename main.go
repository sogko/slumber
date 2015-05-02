package main

import (
	"github.com/sogko/golang-rest-api-server-example/server"
	"github.com/unrolled/render"
)

func main() {

	// initialize server components
	components := server.Components{
		DatabaseSession: server.NewSession(server.DatabaseOptions{
			ServerName:   "localhost",
			DatabaseName: "test-app",
		}),
		Renderer: server.NewRenderer(render.Options{
			IndentJSON: true,
		}),
	}

	// init server and run
	s := server.NewServer(&components)
	s.Run(":3001")

}
