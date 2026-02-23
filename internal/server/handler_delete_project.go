package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerDeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		log.Printf("Project not found")
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting access token: %v", err)
		respondError(w, http.StatusUnauthorized, "cannot fetch access token")
		return
	}
	userId, err := auth.ValidateJWT(accessToken, c.Env.JWT_SECRET)
	if err != nil {
		log.Printf("error validating access token: %v", err)
		respondError(w, http.StatusUnauthorized, "cannot validate access token")
		return
	}
	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		log.Printf("error getting target project from DB: %v", err)
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}

	relation, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	if relation.UserStatus != "admin" {
		log.Printf("you must be a project admin to delete this project")
		respondError(w, http.StatusUnauthorized, "you must be a project admin to delete this project")
		return
	}

	assets, err := c.DB.GetAssetsByProjectID(ctx, project.ID)
	if err != nil {
		log.Printf("error getting assets for project: %s", err)
		respondError(w, http.StatusBadRequest, "could not get project assets for deletion")
		return
	}

// 1. Initialize a slice of strings to hold the keys
	keys := make([]string, 0, len(assets))

	// 2. Loop through the assets to collect their storage paths
	for _, asset := range assets {
		keys = append(keys, asset.StoragePath)
	}

	// 3. Send the entire slice to your S3 client once
	if len(keys) > 0 {
		err = c.S3Client.DeleteObjects(ctx, keys)
		if err != nil {
			log.Printf("error batch deleting assets: %v", err)
			respondError(w, http.StatusConflict, "Error deleting assets, project deletion not complete")
			return
		}
	}
	

	err = c.DB.DeleteProjectByTitle(ctx, projectName)
	if err != nil {
		log.Printf("error deleting project: %v", err)
		respondError(w, http.StatusBadRequest, "error encountered while deleting project")
		return
	}

	respondJSON(w, http.StatusOK, projectName)

}
