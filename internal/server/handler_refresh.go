package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
)

func (c *Config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var refreshReq models.RefreshRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&refreshReq); err != nil {
		log.Printf("could not decode refresh request body: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if refreshReq.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	user, err := c.DB.GetUserFromRefreshToken(ctx, refreshReq.RefreshToken)
	if err != nil {
		log.Printf("invalid refresh token: %v", err)
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, c.Env.JWT_SECRET)
	if err != nil {
		log.Printf("error creating access token from refresh token: %v", err)
		respondError(w, http.StatusConflict, "could not create access token")
		return
	}

	resp := models.RefreshResponse{AccessToken: accessToken}
	respondJSON(w, http.StatusAccepted, resp)
}
