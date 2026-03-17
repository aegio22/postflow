package server

import (
	"net/http"

	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerDeleteProjectMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get authenticated user ID from context (set by middleware)
	userId, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}

	projectName := r.URL.Query().Get("project_name")
	userEmail := r.URL.Query().Get("user_email")
	if projectName == "" {
		logger.Warn(ctx, "project name missing", "operation", "delete_project_member", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	if userEmail == "" {
		logger.Warn(ctx, "user email missing", "operation", "delete_project_member", "user_id", userId.String())
		respondError(w, http.StatusBadRequest, "User not found")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, projectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "delete_project_member", "user_id", userId.String(), "project_name", projectName)
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	removedUser, err := c.DB.GetUserByEmail(ctx, userEmail)
	if err != nil {
		logger.Error(ctx, "user not found", err, "operation", "delete_project_member", "user_id", userId.String(), "target_email", userEmail)
		respondError(w, http.StatusConflict, "could not find user in DB")
		return
	}
	// get both users statuses
	relationReqUser, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "user not in project", err, "operation", "delete_project_member", "user_id", userId.String(), "project_id", project.ID.String())
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}
	relationRmUser, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: removedUser.ID, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "target user not in project", err, "operation", "delete_project_member", "user_id", userId.String(), "target_user_id", removedUser.ID.String())
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}
	//auth checks
	if relationReqUser.UserStatus != "admin" {
		logger.Warn(ctx, "insufficient permissions to remove member", "operation", "delete_project_member", "user_id", userId.String(), "user_status", relationReqUser.UserStatus)
		respondError(w, http.StatusUnauthorized, "you must be a project admin to remove another user")
		return
	}

	if relationRmUser.UserStatus == "admin" && project.CreatedBy != relationReqUser.UserID {
		logger.Warn(ctx, "cannot remove admin without being creator", "operation", "delete_project_member", "user_id", userId.String(), "target_user_id", removedUser.ID.String())
		respondError(w, http.StatusUnauthorized, "you may not remove another admin unless you are the project creator")
		return
	}
	// DB removal once auth checks out
	err = c.DB.RemoveUserFromProject(ctx, database.RemoveUserFromProjectParams{UserID: removedUser.ID, ProjectID: project.ID})
	if err != nil {
		logger.Error(ctx, "failed to remove user from project", err, "operation", "delete_project_member", "user_id", userId.String(), "target_user_id", removedUser.ID.String(), "project_id", project.ID.String())
		respondError(w, http.StatusBadRequest, "user removal from project failed")
		return
	}

	respondJSON(w, http.StatusOK, removedUser.ID)

}
