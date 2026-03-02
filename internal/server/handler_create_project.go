package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerCreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "failed to get user from context")
		return
	}
	var projInfo models.ProjectRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&projInfo)
	if err != nil {
		log.Printf("could not fetch user info from request: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := c.DB.CreateProject(ctx, database.CreateProjectParams{Title: projInfo.Title, Column2: projInfo.Description, CreatedBy: userId})
	if err != nil {
		log.Printf("error creating project: %v", err)
		respondError(w, http.StatusBadRequest, "project creation failed")
		return
	}
	_, err = c.DB.AddNewProjectUser(ctx, database.AddNewProjectUserParams{ProjectID: project.ID, UserID: userId, UserStatus: "admin"})
	if err != nil {
		log.Printf("error setting project author as admin: %v", err)
		respondError(w, http.StatusBadRequest, "project creation failed")
		return
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
