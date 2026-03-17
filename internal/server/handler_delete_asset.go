package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerDeleteAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user ID from context (set by middleware)
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	projectName := r.URL.Query().Get("project_name")
	assetName := r.URL.Query().Get("asset_name")
	if projectName == "" {
		logger.Warn(ctx, "project name missing", "operation", "delete_asset", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	if assetName == "" {
		logger.Warn(ctx, "asset name missing", "operation", "delete_asset", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "Asset not found")
		return
	}
	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "delete_asset", "user_id", userId.String(), "project_name", projectName)
		respondError(w, http.StatusBadRequest, "Project not found")
	}

	relation, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "delete_asset", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	if relation.UserStatus != "admin" && relation.UserStatus != "staff" {
		logger.Warn(ctx, "insufficient permissions to delete asset", "operation", "delete_asset", "user_id", userId.String(), "project_id", project.ID.String(), "user_status", relation.UserStatus)
		respondError(w, http.StatusUnauthorized, "you must be a project admin or staff member to delete this asset")
		return
	}
	asset, err := c.DB.GetAssetByName(ctx, database.GetAssetByNameParams{Name: assetName, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "asset not found", err, "operation", "delete_asset", "user_id", userId.String(), "project_id", project.ID.String(), "asset_name", assetName)
		respondError(w, http.StatusUnauthorized, "could not get asset from DB")
		return
	}

	err = c.S3Client.DeleteObject(ctx, asset.StoragePath)
	if err != nil {
		logger.Error(ctx, "failed to delete asset from S3", err, "operation", "delete_asset", "user_id", userId.String(), "asset_id", asset.ID.String(), "storage_path", asset.StoragePath)
		respondError(w, http.StatusBadRequest, "could not delete asset from s3")
		return
	}

	err = c.DB.DeleteAssetByID(ctx, asset.ID)
	if err != nil {
		logger.Error(ctx, "failed to delete asset from database", err, "operation", "delete_asset", "user_id", userId.String(), "asset_id", asset.ID.String())
		respondError(w, http.StatusBadRequest, "could not delete asset from database")
		return
	}

	logger.Info(ctx, "asset deleted", "operation", "delete_asset", "user_id", userId.String(), "asset_id", asset.ID.String(), "asset_name", assetName)

	respondJSON(w, http.StatusOK, assetName)
}
