package server

import (
	"errors"
	"os"
)

type Env struct {
	DB_URL     string
	JWT_SECRET string
	PORT       string
	AWS_REGION string
	S3_BUCKET  string
}

func LoadEnv() (*Env, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, errors.New("DATABASE_URL not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return nil, errors.New("AWS_REGION not set")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		return nil, errors.New("S3_BUCKET not set")
	}

	return &Env{
		DB_URL:     dbURL,
		JWT_SECRET: jwtSecret,
		PORT:       port,
		AWS_REGION: awsRegion,
		S3_BUCKET:  s3Bucket,
	}, nil
}
