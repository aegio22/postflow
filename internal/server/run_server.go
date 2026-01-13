package server

import "log"

func Run() {
	server, err := CreateServer()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Printf("Starting server on port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)

	}
}
