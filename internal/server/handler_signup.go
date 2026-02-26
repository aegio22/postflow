package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/database"
)

func (c *Config) handlerSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo models.UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("error decoding request: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	//add user to db
	hashedPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		log.Printf("error hashing password for DB storage: %v", err)
		respondError(w, http.StatusConflict, "error hashing password for DB storage")
		return
	}
	newUser, err := c.DB.CreateUser(ctx, database.CreateUserParams{
		Username: userInfo.Username, Email: userInfo.Email, HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("error registering user: %v", err)
		respondError(w, http.StatusBadRequest, "error registering user")
		return
	}
	accessToken, err := auth.MakeJWT(newUser.ID, c.Env.JWT_SECRET)
	if err != nil {
		log.Printf("error making JWT: %v", err)
		respondError(w, http.StatusBadRequest, "error making JWT")
		return
	}
	//create refresh token and add it to the DB
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error making user secret token: %v", err)
		respondError(w, http.StatusBadRequest, "error making user refresh token")
		return
	}
	_, err = c.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{Token: refreshToken, UserID: newUser.ID, ExpiresAt: time.Now().AddDate(0, 0, 60)})
	if err != nil {
		log.Printf("error adding refresh token to database: %v", err)
		respondError(w, http.StatusBadRequest, "error adding refresh token to database")
		return
	}

	respUser := models.DBUserResponse{
		ID:          newUser.ID,
		Username:    newUser.Username,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	//ready and write response
	respondJSON(w, http.StatusCreated, respUser)
}
