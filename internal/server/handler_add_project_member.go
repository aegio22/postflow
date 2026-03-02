package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
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
		log.Printf("could not fetch member info from request: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	project, err := c.DB.GetProjectByTitle(ctx, memberInfo.ProjectName)
	if err != nil {
		log.Printf("error pulling project info from db: %v", err)
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
		log.Printf("error pulling user info from db: %v", err)
		respondError(w, http.StatusBadRequest, "error pulling user info from db")
		return
	}
	_, err = c.DB.AddNewProjectUser(ctx, database.AddNewProjectUserParams{
		ProjectID:  project.ID,
		UserID:     user.ID,
		UserStatus: memberInfo.UserStatus,
	})
	if err != nil {
		log.Printf("error adding project member: %v", err)
		respondError(w, http.StatusBadRequest, "error adding project member")
		return
	}
	responseBody := models.ProjectMemberAddResponse{
		ProjectName: project.Title,
		UserStatus:  memberInfo.UserStatus,
	}

	respondJSON(w, http.StatusCreated, responseBody)
}
