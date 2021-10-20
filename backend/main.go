package main

import (
	"log"
	"memolang/configuration"
	"memolang/server"
)

func main() {

	cfg := configuration.GetConfiguration()

	srv := server.NewServer(cfg)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}