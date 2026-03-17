package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerLsProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "failed to get user from context")
		return
	}
	user, err := c.DB.GetUserByID(ctx, userId)
	if err != nil {
		logger.Error(ctx, "failed to get user info", err, "operation", "ls_projects", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "error pulling user info from db")
		return
	}
	projects, err := c.DB.GetProjectsForUser(ctx, userId)
	if err != nil {
		logger.Error(ctx, "failed to get projects for user", err, "operation", "ls_projects", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "could not get projects from DB for user")
		return
	}
	projectsMap := make(map[string]string)
	for _, proj := range projects {
		status, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: proj.ID})
		if err != nil {
			logger.Error(ctx, "failed to get user status for project", err, "operation", "ls_projects", "user_id", userId.String(), "project_id", proj.ID.String(), "project_title", proj.Title)
			respondError(w, http.StatusConflict, "error getting user status for project")
			return
		}
		projectsMap[proj.Title] = status.UserStatus
	}

	responseBody := models.ProjectsLsResponse{UserName: user.Username, Projects: projectsMap}

	respondJSON(w, http.StatusOK, responseBody)

}
