package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerProjectsUserlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectName := r.URL.Query().Get("project_name")
	if projectName == "" {
		log.Printf("Project not found")
		respondError(w, http.StatusBadRequest, "Project not found")
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
		log.Println("Project not found in database")
		respondError(w, http.StatusBadRequest, "project not found in database")
		return
	}
	_, err = c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}

	users, err := c.DB.GetAllProjectUsers(ctx, project.ID)
	if err != nil {
		log.Printf("error getting project users: %v", err)
		respondError(w, http.StatusConflict, "could not get project users from database")
		return
	}
	userList := make(map[string]string)
	for _, projUser := range users {
		userInfo, err := c.DB.GetUserByID(ctx, projUser.UserID)
		if err != nil {
			log.Printf("error getting user info for user %v: %v", projUser.UserID, err)
			respondError(w, http.StatusConflict, "Error getting user info")
			return
		}
		userList[userInfo.Email] = projUser.UserStatus
	}
	responseBody := models.ProjectsUserlistResponse{Users: userList}
	respondJSON(w, http.StatusOK, responseBody)

}
