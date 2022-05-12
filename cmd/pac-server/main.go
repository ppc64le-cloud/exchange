package main

import (
	"log"

	pacserver "github.com/ppc64le-cloud/exchange/pkg/cmd/pac-server"
)

func main() {
	cmd := pacserver.NewPacServerCommand()
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("failed to run the server: %+v", err)
	}
}
