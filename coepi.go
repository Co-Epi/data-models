package main

import (
	"flag"
	"fmt"

	"github.com/wolkdb/coepi-backend-go/backend"
	"github.com/wolkdb/coepi-backend-go/server"
)

func main() {
	port := flag.Uint("port", uint(server.DefaultPort), "port coepi is listening on")
	project := flag.String("project", backend.DefaultProject, "Google Cloud Project Name")
	instance := flag.String("instance", backend.DefaultInstance, "Google BigTable Instance")

	_, err := server.NewServer(uint16(*port), *project, *instance)
	if err != nil {
		panic(err)
	}
	fmt.Printf("CoEpi Go Server - Listening on port %d...\n", *port)
	for {
	}
}
