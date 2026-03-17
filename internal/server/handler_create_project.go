package server

import (
	"encoding/json"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
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
		logger.Error(ctx, "failed to decode project request", err, "operation", "create_project", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := c.DB.CreateProject(ctx, database.CreateProjectParams{Title: projInfo.Title, Column2: projInfo.Description, CreatedBy: userId})
	if err != nil {
		logger.Error(ctx, "failed to create project", err, "operation", "create_project", "user_id", userId.String(), "project_title", projInfo.Title)
		respondError(w, http.StatusBadRequest, "project creation failed")
		return
	}
	_, err = c.DB.AddNewProjectUser(ctx, database.AddNewProjectUserParams{ProjectID: project.ID, UserID: userId, UserStatus: "admin"})
	if err != nil {
		logger.Error(ctx, "failed to set project author as admin", err, "operation", "create_project", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusBadRequest, "project creation failed")
		return
	}

	logger.Info(ctx, "project created", "operation", "create_project", "user_id", userId.String(), "project_id", project.ID.String(), "project_title", project.Title)

	responseBody := models.ProjectResponse{
		ID:          project.ID,
		Title:       project.Title,
		Description: project.Description.String,
		Status:      project.Status,
		CreatedBy:   project.CreatedBy,
	}
	respondJSON(w, http.StatusCreated, responseBody)

}
