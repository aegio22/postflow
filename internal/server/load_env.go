package server

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_URL     string
	JWT_SECRET string
	PORT       string
}

func LoadEnv() (*Env, error) {
	godotenv.Load() // #nosec G104
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("no jwt secret found")
	}
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("port env missing")
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		return nil, errors.New("db string missing")
	}

	return &Env{
		PORT:   port,
		DB_URL: db,
	}, nil

}
