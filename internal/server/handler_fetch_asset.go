package server

import (
	"log"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerViewAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assetName := r.URL.Query().Get("asset_name")
	projectName := r.URL.Query().Get("project_name")

	if assetName == "" || projectName == "" {
		respondError(w, http.StatusBadRequest, "missing asset_name or project_name")
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting access token from request: %v", err)
		respondError(w, http.StatusUnauthorized, "could not find token")
		return
	}

	userId, err := auth.ValidateJWT(accessToken, c.Env.JWT_SECRET)
	if err != nil {
		log.Printf("error validating access token: %v", err)
		respondError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		log.Printf("error getting project from title: %s", err)
		respondError(w, http.StatusBadRequest, "could not get project from provided title")
		return
	}

	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{
		UserID: userId, ProjectID: project.ID,
	})
	if err != nil {
		log.Printf("error finding user info in project: %s", err)
		respondError(w, http.StatusUnauthorized, "user not found in project")
		return
	}

	asset, err := c.DB.GetAssetByName(ctx, database.GetAssetByNameParams{
		Name:      assetName,
		ProjectID: project.ID,
	})
	if err != nil {
		log.Printf("error fetching asset from the database: %v", err)
		respondError(w, http.StatusBadRequest, "requested asset not found in database")
		return
	}

	// Generate presigned S3 GET URL for this asset
	const ttl = 15 * time.Minute
	downloadURL, err := c.S3Client.PresignDownload(ctx, asset.StoragePath, ttl)
	if err != nil {
		log.Printf("failed to generate presigned download URL: %v", err)
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
