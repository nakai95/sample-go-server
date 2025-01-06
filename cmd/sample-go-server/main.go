package main

import (
	"flag"
	"net"

	"sample-go-server/internal/infrastructure/server"
)

func main() {
	port := flag.String("port", "8080", "Port for test HTTP server")
	flag.Parse()

	e := server.NewRouter()

	// And we serve HTTP until the world ends.
	e.Logger.Fatal(e.Start(net.JoinHostPort("0.0.0.0", *port)))
}
