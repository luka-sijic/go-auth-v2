package main

import (
	"app/internal/server"
	"log"
)

func main() {
	log.Println("users table is ready on all shards")
	server.Start()
}
