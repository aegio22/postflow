package server

import (
	"encoding/json"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var refreshReq models.RefreshRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&refreshReq); err != nil {
		logger.Error(ctx, "failed to decode refresh request", err, "operation", "refresh")
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if refreshReq.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	user, err := c.DB.GetUserFromRefreshToken(ctx, refreshReq.RefreshToken)
	if err != nil {
		logger.Error(ctx, "invalid refresh token", err, "operation", "refresh")
		respondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, c.Env.JWT_SECRET)
	if err != nil {
		logger.Error(ctx, "failed to create access token", err, "operation", "refresh", "user_id", user.ID.String())
		respondError(w, http.StatusConflict, "could not create access token")
		return
	}

	resp := models.RefreshResponse{AccessToken: accessToken}
	respondJSON(w, http.StatusAccepted, resp)
}
