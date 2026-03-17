package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerProjectsUserlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user
	authenticatedUserID, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		respondError(w, http.StatusBadRequest, "project_name required")
		return
	}

	// Fix: Actually fetch the project!
	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		respondError(w, http.StatusNotFound, "project not found")
		return
	}

	// Security: Verify authenticated user is a member
	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{
		UserID:    authenticatedUserID,
		ProjectID: project.ID,
	})
	if err != nil {
		respondError(w, http.StatusForbidden, "not a member of this project")
		return
	}

	users, err := c.DB.GetAllProjectUsers(ctx, project.ID)
	if err != nil {
		logger.Error(ctx, "failed to get project users", err, "operation", "projects_userlist", "user_id", authenticatedUserID.String(), "project_id", project.ID.String())
		respondError(w, http.StatusConflict, "could not get project users from database")
		return
	}
	userList := make(map[string]string)
	for _, projUser := range users {
		userInfo, err := c.DB.GetUserByID(ctx, projUser.UserID)
		if err != nil {
			logger.Error(ctx, "failed to get user info", err, "operation", "projects_userlist", "user_id", authenticatedUserID.String(), "target_user_id", projUser.UserID.String())
			respondError(w, http.StatusConflict, "Error getting user info")
			return
		}
		userList[userInfo.Email] = projUser.UserStatus
	}
	responseBody := models.ProjectsUserlistResponse{Users: userList}
	respondJSON(w, http.StatusOK, responseBody)

}
