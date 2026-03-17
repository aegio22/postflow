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

func (c *Config) handlerSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo models.UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		logger.Error(ctx, "failed to decode signup request", err, "operation", "signup")
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	//add user to db
	hashedPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		logger.Error(ctx, "failed to hash password", err, "operation", "signup", "email", userInfo.Email)
		respondError(w, http.StatusConflict, "error hashing password for DB storage")
		return
	}
	newUser, err := c.DB.CreateUser(ctx, database.CreateUserParams{
		Username: userInfo.Username, Email: userInfo.Email, HashedPassword: hashedPassword,
	})
	if err != nil {
		logger.Error(ctx, "failed to create user", err, "operation", "signup", "email", userInfo.Email, "username", userInfo.Username)
		respondError(w, http.StatusBadRequest, "error registering user")
		return
	}
	accessToken, err := auth.MakeJWT(newUser.ID, c.Env.JWT_SECRET)
	if err != nil {
		logger.Error(ctx, "failed to create JWT", err, "operation", "signup", "user_id", newUser.ID.String())
		respondError(w, http.StatusBadRequest, "error making JWT")
		return
	}
	//create refresh token and add it to the DB
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		logger.Error(ctx, "failed to create refresh token", err, "operation", "signup", "user_id", newUser.ID.String())
		respondError(w, http.StatusBadRequest, "error making user refresh token")
		return
	}
	_, err = c.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{Token: refreshToken, UserID: newUser.ID, ExpiresAt: time.Now().AddDate(0, 0, 60)})
	if err != nil {
		logger.Error(ctx, "failed to store refresh token", err, "operation", "signup", "user_id", newUser.ID.String())
		respondError(w, http.StatusBadRequest, "error adding refresh token to database")
		return
	}

	respUser := models.DBUserResponse{
		ID:           newUser.ID,
		Username:     newUser.Username,
		CreatedAt:    newUser.CreatedAt,
		UpdatedAt:    newUser.UpdatedAt,
		Email:        newUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	logger.Info(ctx, "user signup successful", "operation", "signup", "user_id", newUser.ID.String(), "email", newUser.Email)
	//ready and write response
	respondJSON(w, http.StatusCreated, respUser)
}
