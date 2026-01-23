package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerDeleteAsset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectName := r.URL.Query().Get("project_name")
	assetName := r.URL.Query().Get("asset_name")
	if projectName == "" {
		log.Printf("Project not found")
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	if assetName == "" {
		log.Printf("Asset not found")
		respondError(w, http.StatusBadRequest, "Asset not found")
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

	relation, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	if relation.UserStatus != "admin" && relation.UserStatus != "staff" {
		log.Printf("you must be a project admin or staff member to delete this asset")
		respondError(w, http.StatusUnauthorized, "you must be a project admin or staff member to delete this asset")
		return
	}
	asset, err := c.DB.GetAssetByName(ctx, database.GetAssetByNameParams{Name: assetName, ProjectID: project.ID})
	if err != nil {
		log.Printf("error getting asset from DB: %s", err)
		respondError(w, http.StatusUnauthorized, "could not get asset from DB")
		return
	}

	err = c.S3Client.DeleteObject(ctx, asset.StoragePath)
	if err != nil {
		log.Printf("error deleting asset from s3: %s", err)
		respondError(w, http.StatusBadRequest, "could not delete asset from s3")
		return
	}

	err = c.DB.DeleteAssetByID(ctx, asset.ID)
	if err != nil {
		log.Printf("error deleting asset from database: %s", err)
		respondError(w, http.StatusBadRequest, "could not delete asset from database")
		return
	}
	respondJSON(w, http.StatusOK, assetName)
}
