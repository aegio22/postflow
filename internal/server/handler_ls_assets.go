package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerLsAssets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user ID from context (set by middleware)
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		logger.Warn(ctx, "project name missing in request", "operation", "ls_assets", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "ls_assets", "user_id", userId.String(), "project_name", projectName)
		respondError(w, http.StatusBadRequest, "Project not found")
	}

	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "ls_assets", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	assets, err := c.DB.GetAssetsByProjectName(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "failed to get assets", err, "operation", "ls_assets", "user_id", userId.String(), "project_name", projectName)
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
