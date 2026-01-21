package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerUploadAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var assetInfo models.AssetRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&assetInfo)
	if err != nil {
		log.Printf("could not fetch asset info from request: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
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

	projectId, err := c.DB.GetProjectByTitle(ctx, assetInfo.ProjectName)
	if err != nil {
		log.Printf("error getting project id from title: %s", err)
		respondError(w, http.StatusBadRequest, "project not found in database")
		return
	}
	usersProjects, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: projectId.ID})
	if err != nil {
		log.Printf("error getting user project relation: %s", err)
		respondError(w, http.StatusBadRequest, "user project relation not found")
		return
	}
	if usersProjects.UserStatus != "admin" && usersProjects.UserStatus != "staff" {
		log.Println("must be a staff or admin user to upload to this project")
		respondError(w, http.StatusUnauthorized, "must be a staff or admin user to upload to this project")
		return
	}
	asset, err := c.DB.CreateAsset(ctx, database.CreateAssetParams{
		ProjectID:   projectId.ID,
		Name:        assetInfo.AssetName,
		StoragePath: "",
		Tags:        assetInfo.Tag,
		CreatedBy:   userId,
	})
	if err != nil {
		log.Printf("error adding asset to db: %v", err)
		respondError(w, http.StatusConflict, "error adding asset to database")
		return
	}

	//Generate S3 key (where file will be stored)
	s3Key := fmt.Sprintf("projects/%s/assets/%s/%s",
		asset.ProjectID.String(),
		asset.ID.String(),
		assetInfo.Filepath,
	)

	//Generate presigned upload URL (client will PUT file here)
	uploadURL, err := c.S3Client.PresignUpload(ctx, s3Key, 15*time.Minute)
	if err != nil {
		log.Printf("failed to generate upload URL: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to generate upload URL")
		return
	}
	err = c.DB.UpdateAssetStoragePath(ctx, database.UpdateAssetStoragePathParams{ID: asset.ID, StoragePath: s3Key})
	if err != nil {
		log.Printf("failed to update storage path: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to update storage path")
	}
	responseBody := models.AssetResponse{
		AssetID:   asset.ID.String(),
		UploadURL: uploadURL,
		S3Key:     s3Key,
		ExpiresIn: 900, // 15 minutes
	}
	respondJSON(w, http.StatusCreated, responseBody)

}
