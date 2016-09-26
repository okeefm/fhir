package main

import "github.com/intervention-engine/fhir/server"

func main() {
	s := server.NewServer("localhost")

	s.Run(server.DefaultConfig)
}
