package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerLsProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	user, err := c.DB.GetUserByID(ctx, userId)
	if err != nil {
		log.Printf("error pulling user info from db: %v", err)
		respondError(w, http.StatusBadRequest, "error pulling user info from db")
		return
	}
	projects, err := c.DB.GetProjectsForUser(ctx, userId)
	if err != nil {
		log.Printf("error fetching projects from DB: %v", err)
		respondError(w, http.StatusBadRequest, "could not get projects from DB for user")
		return
	}
	projectsMap := make(map[string]string)
	for _, proj := range projects {
		status, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: proj.ID})
		if err != nil {
			log.Printf("could not get user status for project %v", proj.Title)
			respondError(w, http.StatusConflict, "error getting user status for project")
			return
		}
		projectsMap[proj.Title] = status.UserStatus
	}

	responseBody := models.ProjectsLsResponse{UserName: user.Username, Projects: projectsMap}

	respondJSON(w, http.StatusOK, responseBody)

}
