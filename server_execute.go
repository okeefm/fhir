package main

import "github.com/okeefm/fhir/server"

func main() {
	s := server.NewServer("localhost")

	s.Run(server.DefaultConfig)
}
