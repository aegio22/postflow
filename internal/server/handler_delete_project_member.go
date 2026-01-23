package server

import (
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerDeleteProjectMember(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectName := r.URL.Query().Get("project_name")
	userEmail := r.URL.Query().Get("user_email")
	if projectName == "" {
		log.Printf("Project not found")
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	if userEmail == "" {
		log.Printf("User not found")
		respondError(w, http.StatusBadRequest, "User not found")
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
		log.Printf("error getting target project from DB: %v", err)
		respondError(w, http.StatusBadRequest, "Project not found")
		return
	}
	removedUser, err := c.DB.GetUserByEmail(ctx, userEmail)
	if err != nil {
		log.Printf("error getting user for removal from DB: %s", err)
		respondError(w, http.StatusConflict, "could not find user in DB")
		return
	}
	// get both users statuses
	relationReqUser, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: userId, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}
	relationRmUser, err := c.DB.GetUserProjectRelation(ctx, database.GetUserProjectRelationParams{UserID: removedUser.ID, ProjectID: project.ID})
	if err != nil {
		log.Printf("error finding user project relation: %v", err)
		respondError(w, http.StatusUnauthorized, "user is not a project member")
		return
	}
	//auth checks
	if relationReqUser.UserStatus != "admin" {
		log.Printf("you must be a project admin to remove another user")
		respondError(w, http.StatusUnauthorized, "you must be a project admin to remove another user")
		return
	}

	if relationRmUser.UserStatus == "admin" && project.CreatedBy != relationReqUser.UserID {
		log.Printf("you may not remove another admin unless you are the project creator")
		respondError(w, http.StatusUnauthorized, "you may not remove another admin unless you are the project creator")
		return
	}
	// DB removal once auth checks out
	err = c.DB.RemoveUserFromProject(ctx, database.RemoveUserFromProjectParams{UserID: removedUser.ID, ProjectID: project.ID})
	if err != nil {
		log.Printf("error removing user from project: %v", err)
		respondError(w, http.StatusBadRequest, "user removal from project failed")
		return
	}

	respondJSON(w, http.StatusOK, removedUser.ID)

}
