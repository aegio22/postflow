package main

import (
	"log"

	"github.com/aegio22/postflow/internal/server"
)

func main() {
	server, err := server.CreateServer()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Printf("Starting server on port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
