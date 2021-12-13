package main

import (
	"context"
	firebase "firebase.google.com/go"
	"log"
	"memolang/configuration"
	"memolang/server"
)

func main() {

	cfg := configuration.GetConfiguration()

	firebaseConfig := &firebase.Config{ProjectID: "memolang-97db1"}

	app, err := firebase.NewApp(context.Background(), firebaseConfig)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	srv := server.NewServer(cfg, app)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
