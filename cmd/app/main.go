package main

import (
	"app/internal/database"
	"app/internal/server"
	"app/pkg/auto"
	"log"
)

func main() {
	database.Connect()
	defer database.Close()

	if err := auto.Init(); err != nil {
		log.Fatalf("could not create users table: %v", err)
	}
	log.Println("users table is ready on all shards")
	server.Start()
}
