package server

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DB_URL                string
	JWT_SECRET            string
	BASE_URL              string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	S3_BUCKET             string
	AWS_REGION            string
}

func LoadEnv() (*Env, error) {
	godotenv.Load() // #nosec G104
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("no jwt secret found")
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	if awsAccessKey == "" {
		return nil, errors.New("No AWS access key found")
	}

	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsSecretKey == "" {
		return nil, errors.New("No AWS secret key found")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		return nil, errors.New("No s3 bucket found")
	}

	awsRegion := os.Getenv("AWS_REGION")
	if awsRegion == "" {
		return nil, errors.New("No AWS region found")
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
		BASE_URL:              baseURL,
		DB_URL:                db,
		JWT_SECRET:            jwtSecret,
		AWS_ACCESS_KEY_ID:     awsAccessKey,
		AWS_SECRET_ACCESS_KEY: awsSecretKey,
		S3_BUCKET:             s3Bucket,
		AWS_REGION:            awsRegion,
	}, nil

}
