package server

import (
	"net/http"
	"time"
)

func CreateServer() (*http.Server, error) {
	cfg, err := CreateConfig()
	if err != nil {
		return nil, err
	}

	r := cfg.NewRouter()

	//initialize core endpoints and handlers
	r.HandleFunc("POST /api/signup", cfg.handlerSignUp)
	//initialize and start server
	server := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return server, nil
}
