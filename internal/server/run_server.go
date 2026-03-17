package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/aegio22/postflow/internal/logger"
	"github.com/aegio22/postflow/internal/storage"
	_ "github.com/lib/pq" // Required for the postgres driver
	"github.com/pressly/goose/v3"
)

func Run(args []string) error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	// Initialize logger for startup
	log := logger.NewLogger()
	ctx := context.Background()

	// 1. Run migrations before starting the server logic
	if err := applyMigrations(ctx, log, dbURL); err != nil {
		// We return the error so RunCLI can print it properly
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	// 2. Initialize and start the server
	server, err := CreateServer()
	if err != nil {
		return fmt.Errorf("failed to create server: %w", err)
	}

	log.Info("starting server", "addr", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		return fmt.Errorf("server crashed: %w", err)
	}

	return nil
}

func applyMigrations(ctx context.Context, log *slog.Logger, dbURL string) error {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	// Tell goose to use the filesystem embedded in our binary
	goose.SetBaseFS(storage.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	log.InfoContext(ctx, "checking database schema")

	if err := goose.Up(db, "sql/schema"); err != nil {
		return err
	}

	log.InfoContext(ctx, "database ready")
	return nil
}
