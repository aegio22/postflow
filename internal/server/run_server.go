package server

import "log"

// included args simply to fit the cliCommand interface
func Run(args []string) error {
	server, err := CreateServer()
	if err != nil {
		log.Fatalf("Error: %s", err)
		return err
	}

	log.Printf("Starting server on port %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return err
	}
	return nil
}
