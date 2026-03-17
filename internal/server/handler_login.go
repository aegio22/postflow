package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
	"github.com/aegio22/postflow/internal/logger"
)

func (c *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo models.UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		logger.Error(ctx, "failed to decode login request", err, "operation", "login")
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := c.DB.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		logger.Error(ctx, "user not found", err, "operation", "login", "email", userInfo.Email)
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, c.Env.JWT_SECRET)
	if err != nil {
		logger.Error(ctx, "failed to create access token", err, "operation", "login", "user_id", user.ID.String())
		respondError(w, http.StatusConflict, "could not create access token")
		return
	}
	passwordMatch, err := auth.CheckPasswordHash(userInfo.Password, user.HashedPassword)
	if err != nil {
		logger.Error(ctx, "password verification failed", err, "operation", "login", "user_id", user.ID.String())
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if passwordMatch != true {
		respondError(w, http.StatusUnauthorized, "password does not match our records")
		return
	}

	refreshRecord, err := c.DB.GetTokenFromUserID(ctx, user.ID)
	if err != nil {
		logger.Error(ctx, "no refresh token found", err, "operation", "login", "user_id", user.ID.String())
		respondError(w, http.StatusUnauthorized, "no refresh token found in DB for given user")
		return
	}
	//Because of how sql.NullStrings work, RevokedAt is valid if the token HAS BEEN revoked
	if refreshRecord.RevokedAt.Valid {
		logger.Warn(ctx, "refresh token is revoked", "operation", "login", "user_id", user.ID.String())
		respondError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}
	activeRefreshToken := refreshRecord.Token
	if refreshRecord.ExpiresAt.Before(time.Now()) {
		err = c.DB.RevokeToken(ctx, refreshRecord.Token)
		if err != nil {
			logger.Error(ctx, "failed to revoke expired token", err, "operation", "login", "user_id", user.ID.String())
			respondError(w, http.StatusConflict, "error revoking expired refresh token")
			return
		}
		newRefreshToken, err := auth.MakeRefreshToken()
		if err != nil {
			logger.Error(ctx, "failed to create new refresh token", err, "operation", "login", "user_id", user.ID.String())
			respondError(w, http.StatusConflict, "old refresh token expired, and there is an error making a new one")
			return
		}
		_, err = c.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{Token: newRefreshToken, UserID: user.ID, ExpiresAt: time.Now().AddDate(0, 0, 60)})
		if err != nil {
			logger.Error(ctx, "failed to store new refresh token", err, "operation", "login", "user_id", user.ID.String())
			respondError(w, http.StatusBadRequest, "error adding refresh token to database")
			return
		}
		activeRefreshToken = newRefreshToken
	}

	logger.Info(ctx, "user login successful", "operation", "login", "user_id", user.ID.String())

	var loginResponse models.LoginResponse
	loginResponse.AccessToken = accessToken
	loginResponse.RefreshToken = activeRefreshToken
	respondJSON(w, http.StatusAccepted, loginResponse)

}
