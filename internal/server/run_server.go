package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"io/fs"
	"github.com/aegio22/postflow" // This pulls in the root embed.go
	"github.com/pressly/goose/v3"
	_ "github.com/lib/pq" // Required for the postgres driver
)

func Run(args []string) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// 1. Run migrations before starting the server logic
	if err := applyMigrations(dbURL); err != nil {
		// We return the error so RunCLI can print it properly
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// 2. Initialize and start the server
	server, err := CreateServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	log.Printf("Starting PostFlow server on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("server crashed: %w", err)
	}
	
	return nil
}

func applyMigrations(dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Tell goose to use the filesystem embedded in our binary
	goose.SetBaseFS(postflow.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	log.Println("Checking database schema...")
	// Wrap the FS so "sql/schema" is the root
	subFS, err := fs.Sub(postflow.MigrationsFS, "sql/schema")
	if err != nil {
		return err
	}
	if err := goose.Up(db, "."); err != nil { // Now use "."
		return err
	}
	
	log.Println("Database is ready.")
	return nil
}
