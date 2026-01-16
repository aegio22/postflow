package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerCreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var projInfo models.ProjectRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&projInfo)
	if err != nil {
		log.Printf("could not fetch user info from request: %v", err)
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
	project, err := c.DB.CreateProject(ctx, database.CreateProjectParams{Title: projInfo.Title, Column2: projInfo.Description, CreatedBy: userId})
	if err != nil {
		log.Printf("error creating project: %v", err)
		respondError(w, http.StatusBadRequest, "project creation failed")
	}
	responseBody := models.ProjectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description.String,
		Status:      project.Status,
		CreatedBy:   project.CreatedBy,
	}
	respondJSON(w, http.StatusCreated, responseBody)

}
