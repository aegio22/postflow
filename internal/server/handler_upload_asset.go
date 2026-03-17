package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
	"github.com/google/uuid"
)

func (c *Config) handlerUploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user ID from context (set by middleware)
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	var assetInfo models.AssetRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&assetInfo)
	if err != nil {
		logger.Error(ctx, "failed to decode asset upload request", err, "operation", "upload_asset", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	projectId, err := c.DB.GetProjectByTitle(ctx, assetInfo.ProjectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "upload_asset", "user_id", userId.String(), "project_name", assetInfo.ProjectName)
		respondError(w, http.StatusBadRequest, "project not found in database")
		return
	}
	usersProjects, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: projectId.ID})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "upload_asset", "user_id", userId.String(), "project_id", projectId.ID.String())
		respondError(w, http.StatusBadRequest, "user project relation not found")
		return
	}
	if usersProjects.UserStatus != "admin" && usersProjects.UserStatus != "staff" {
		logger.Warn(ctx, "insufficient permissions for asset upload", "operation", "upload_asset", "user_id", userId.String(), "project_id", projectId.ID.String(), "user_status", usersProjects.UserStatus)
		respondError(w, http.StatusUnauthorized, "must be a staff or admin user to upload to this project")
		return
	}
	resp, err := c.createAssetWithUploadURL(
		ctx,
		projectId.ID,
		userId,
		assetInfo.AssetName,
		assetInfo.Tag,
		assetInfo.Filepath,
	)
	if err != nil {
		logger.Error(ctx, "failed to create asset", err, "operation", "upload_asset", "user_id", userId.String(), "project_id", projectId.ID.String(), "asset_name", assetInfo.AssetName)
		respondError(w, http.StatusConflict, "error adding asset to database")
		return
	}

	logger.Info(ctx, "asset upload URL created", "operation", "upload_asset", "user_id", userId.String(), "project_id", projectId.ID.String(), "asset_name", assetInfo.AssetName)

	respondJSON(w, http.StatusCreated, resp)

}

func (c *Config) createAssetWithUploadURL(
	ctx context.Context,
	projectID uuid.UUID,
	userID uuid.UUID,
	assetName string,
	tag string,
	filepath string,
) (models.AssetResponse, error) {
	// 1) Insert asset row
	asset, err := c.DB.CreateAsset(ctx, database.CreateAssetParams{
		ProjectID:   projectID,
		Name:        assetName,
		StoragePath: "",
		Tags:        tag,
		CreatedBy:   userID,
	})
	if err != nil {
		return models.AssetResponse{}, fmt.Errorf("create asset: %w", err)
	}

	s3Key := fmt.Sprintf("projects/%s/assets/%s/%s",
		asset.ProjectID.String(),
		asset.ID.String(),
		filepath,
	)

	const ttl = 15 * time.Minute
	uploadURL, err := c.S3Client.PresignUpload(ctx, s3Key, ttl)
	if err != nil {
		return models.AssetResponse{}, fmt.Errorf("presign upload: %w", err)
	}

	if err := c.DB.UpdateAssetStoragePath(ctx, database.UpdateAssetStoragePathParams{
		ID:          asset.ID,
		StoragePath: s3Key,
	}); err != nil {
		return models.AssetResponse{}, fmt.Errorf("update storage path: %w", err)
	}

	resp := models.AssetResponse{
		AssetID:   asset.ID.String(),
		UploadURL: uploadURL,
		S3Key:     s3Key,
		ExpiresIn: int(ttl.Seconds()),
	}
	return resp, nil
}
