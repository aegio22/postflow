package server

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_URL     string
	JWT_SECRET string
	BASE_URL   string
}

func LoadEnv() (*Env, error) {
	godotenv.Load() // #nosec G104
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("no jwt secret found")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return nil, errors.New("baseurl env missing")
	}

	db := os.Getenv("DATABASE_URL")
	if db == "" {
		return nil, errors.New("db string missing")
	}

	return &Env{
		BASE_URL: baseURL,
		DB_URL:   db,
	}, nil

}
