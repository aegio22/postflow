package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/routes"
)

func CreateServer() (*http.Server, error) {
	cfg, err := CreateConfig()
	if err != nil {
		return nil, err
	}

	r := cfg.NewRouter()

	//initialize core endpoints and handlers
	r.HandleFunc("POST "+routes.SignUp, cfg.handlerSignUp)
	r.HandleFunc("POST "+routes.Login, cfg.handlerLogin)
	r.HandleFunc("POST "+routes.Projects, cfg.handlerCreateProject)
	r.HandleFunc("POST "+routes.ProjectMembers, cfg.handlerAddProjectMember)
	r.HandleFunc("POST "+routes.Assets, cfg.handlerUploadAsset)
	r.HandleFunc("GET "+routes.ViewAssets, cfg.handlerViewAsset)
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

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
