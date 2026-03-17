package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerDeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		logger.Warn(ctx, "project name missing in delete request", "operation", "delete_project")
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}

	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "failed to get user from context")
		return
	}
	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "delete_project", "user_id", userId.String(), "project_name", projectName)
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}

	relation, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	if relation.UserStatus != "admin" {
		logger.Warn(ctx, "insufficient permissions to delete project", "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String(), "user_status", relation.UserStatus)
		respondError(w, http.StatusUnauthorized, "you must be a project admin to delete this project")
		return
	}

	assets, err := c.DB.GetAssetsByProjectID(ctx, project.ID)
	if err != nil {
		logger.Error(ctx, "failed to get assets for project", err, "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String())
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
			logger.Error(ctx, "failed to delete assets from S3", err, "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String(), "asset_count", len(keys))
			respondError(w, http.StatusConflict, "Error deleting assets, project deletion not complete")
			return
		}
	}

	err = c.DB.DeleteProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "failed to delete project from database", err, "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusBadRequest, "error encountered while deleting project")
		return
	}

	logger.Info(ctx, "project deleted", "operation", "delete_project", "user_id", userId.String(), "project_id", project.ID.String(), "project_name", projectName)

	respondJSON(w, http.StatusOK, projectName)

}
