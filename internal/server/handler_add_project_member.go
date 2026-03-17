package server

import (
	"encoding/json"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerAddProjectMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticatedUserID, ok := getUserID(ctx)
	if !ok {
		respondError(w, http.StatusInternalServerError, "authentication error")
		return
	}
	var memberInfo models.AddUserRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&memberInfo)
	if err != nil {
		logger.Error(ctx, "failed to decode add member request", err, "operation", "add_project_member", "user_id", authenticatedUserID.String())
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, memberInfo.ProjectName)
	if err != nil {
		logger.Error(ctx, "project not found", err, "operation", "add_project_member", "user_id", authenticatedUserID.String(), "project_name", memberInfo.ProjectName)
		respondError(w, http.StatusBadRequest, "error pulling project info from db")
		return
	}
	relation, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{
		UserID:    authenticatedUserID,
		ProjectID: project.ID,
	})
	if err != nil || (relation.UserStatus != "admin" && relation.UserStatus != "staff") {
		respondError(w, http.StatusForbidden, "only admins/staff can add members")
		return
	}
	user, err := c.DB.GetUserByEmail(ctx, memberInfo.UserEmail)
	if err != nil {
		logger.Error(ctx, "user not found", err, "operation", "add_project_member", "user_id", authenticatedUserID.String(), "target_email", memberInfo.UserEmail)
		respondError(w, http.StatusBadRequest, "error pulling user info from db")
		return
	}
	_, err = c.DB.AddNewProjectUser(ctx, database.AddNewProjectUserParams{
		ProjectID:  project.ID,
		UserID:     user.ID,
		UserStatus: memberInfo.UserStatus,
	})
	if err != nil {
		logger.Error(ctx, "failed to add project member", err, "operation", "add_project_member", "user_id", authenticatedUserID.String(), "project_id", project.ID.String(), "target_user_id", user.ID.String())
		respondError(w, http.StatusBadRequest, "error adding project member")
		return
	}
	responseBody := models.ProjectMemberAddResponse{
		ProjectName: project.Title,
		UserStatus:  memberInfo.UserStatus,
	}

	respondJSON(w, http.StatusCreated, responseBody)
}
