package main

import (
	. "github.com/sogko/rest-api-server/server"
	"github.com/unrolled/render"
)

func main() {

	// initialize server components
	components := ServerComponents{
		DatabaseSession: NewSession(DatabaseOptions{
			ServerName:   "localhost",
			DatabaseName: "test-app",
		}),
		Renderer: NewRenderer(render.Options{
			IndentJSON: true,
		}),
	}

	// init server and run
	server := NewServer(&components)
	server.Run(":3001")

}
