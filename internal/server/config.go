package server

import (
	"context"
	"database/sql"

	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/storage"
	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/lib/pq"
)

type Config struct {
	DB       *database.Queries
	Env      *Env
	S3Client *storage.S3Client
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

	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(env.AWS_REGION),
	)
	if err != nil {
		return nil, err
	}

	s3Client := storage.NewS3(awsCfg, env.S3_BUCKET)

	return &Config{
		DB:       database.New(db),
		Env:      env,
		S3Client: s3Client,
	}, nil
}
