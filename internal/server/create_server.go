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

	// Public endpoints (no auth required)
	r.HandleFunc("POST "+routes.SignUp, cfg.handlerSignUp)
	r.HandleFunc("POST "+routes.Login, cfg.handlerLogin)
	r.HandleFunc("POST "+routes.Refresh, cfg.handlerRefresh)

	// Protected endpoints (auth required)
	r.HandleFunc("POST "+routes.Projects, cfg.requireAuth(cfg.handlerCreateProject))
	r.HandleFunc("GET "+routes.Projects, cfg.requireAuth(cfg.handlerLsProjects))
	r.HandleFunc("DELETE "+routes.Projects, cfg.requireAuth(cfg.handlerDeleteProject))

	r.HandleFunc("POST "+routes.ProjectMembers, cfg.requireAuth(cfg.handlerAddProjectMember))
	r.HandleFunc("GET "+routes.ProjectMembers, cfg.requireAuth(cfg.handlerProjectsUserlist))
	r.HandleFunc("DELETE "+routes.ProjectMembers, cfg.requireAuth(cfg.handlerDeleteProjectMember))

	r.HandleFunc("POST "+routes.Assets, cfg.requireAuth(cfg.handlerUploadAsset))
	r.HandleFunc("GET "+routes.Assets, cfg.requireAuth(cfg.handlerLsAssets))
	r.HandleFunc("DELETE "+routes.Assets, cfg.requireAuth(cfg.handlerDeleteAsset))
	r.HandleFunc("GET "+routes.ViewAssets, cfg.requireAuth(cfg.handlerViewAsset))

	server := &http.Server{
		Addr:              cfg.Env.PORT,
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
