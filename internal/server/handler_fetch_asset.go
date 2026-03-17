package server

import (
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerViewAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user ID from context (set by middleware)
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	assetName := r.URL.Query().Get("asset_name")
	projectName := r.URL.Query().Get("project_name")

	if assetName == "" || projectName == "" {
		respondError(w, http.StatusBadRequest, "missing asset_name or project_name")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "view_asset", "user_id", userId.String(), "project_name", projectName)
		respondError(w, http.StatusBadRequest, "could not get project from provided title")
		return
	}

	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{
		UserID: userId, ProjectID: project.ID,
	})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "view_asset", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusUnauthorized, "user not found in project")
		return
	}

	asset, err := c.DB.GetAssetByName(ctx, database.GetAssetByNameParams{
		Name:      assetName,
		ProjectID: project.ID,
	})
	if err != nil {
		logger.Error(ctx, "asset not found", err, "operation", "view_asset", "user_id", userId.String(), "project_id", project.ID.String(), "asset_name", assetName)
		respondError(w, http.StatusBadRequest, "requested asset not found in database")
		return
	}

	// Generate presigned S3 GET URL for this asset
	const ttl = 15 * time.Minute
	downloadURL, err := c.S3Client.PresignDownload(ctx, asset.StoragePath, ttl, assetName)
	if err != nil {
		logger.Error(ctx, "failed to generate presigned URL", err, "operation", "view_asset", "user_id", userId.String(), "asset_id", asset.ID.String())
		respondError(w, http.StatusInternalServerError, "failed to generate asset URL")
		return
	}

	resp := models.AssetResponse{
		AssetID:   asset.ID.String(),
		UploadURL: downloadURL,       // reuse field as "URL" for viewing/downloading
		S3Key:     asset.StoragePath, // optional, but nice to return
		ExpiresIn: int(ttl.Seconds()),
	}

	respondJSON(w, http.StatusOK, resp)
}
