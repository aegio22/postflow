package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerLsAssets(w http.ResponseWriter, r *http.Request) {
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
	}

	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	assets, err := c.DB.GetAssetsByProjectName(ctx, projectName)
	if err != nil {
		log.Printf("error fetching assets from DB: %v", err)
		respondError(w, http.StatusBadRequest, "could not get assets from DB for project")
		return
	}

	assetsMap := make(map[string]string)
	for _, asset := range assets {
		assetName := asset.Name
		tag := asset.Tags
		assetsMap[assetName] = tag
	}

	responseBody := models.AssetsLsResponse{Assets: assetsMap}

	respondJSON(w, http.StatusOK, responseBody)

}
