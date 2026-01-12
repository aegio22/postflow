package server

import (
	"database/sql"

	"github.com/aegio22/postflow/internal/database"
)

type Config struct {
	DB  *database.Queries
	Env *Env
}

func CreateConfig() (*Config, error) {
	env, err := LoadEnv()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", env.DB_URL)
	if err != nil {
		return nil, err
	}

	dbQueries := database.New(db)
	return &Config{
		DB:  dbQueries,
		Env: env,
	}, nil
}
