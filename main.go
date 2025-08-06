package main

import (
	"log"

	"doh-server/config"
	"doh-server/server"
)

func main() {
	cfg := config.Load()
	handler := server.NewRouter(cfg)
	log.Fatal(server.Start(cfg, handler))
}
