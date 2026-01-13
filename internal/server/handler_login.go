package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
)

func (c *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("could not fetch user info from request: %v", err)
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := c.DB.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Printf("error getting user: %v", err)
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	accessToken, err := auth.MakeJWT(user.ID, c.Env.JWT_SECRET)
	if err != nil {
		log.Printf("error creating access token: %v", err)
		respondError(w, http.StatusConflict, "could not create access token")
		return
	}
	passwordMatch, err := auth.CheckPasswordHash(userInfo.Password, user.HashedPassword)
	if err != nil {
		log.Printf("error comparing password with hash: %v", err)
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}
	if passwordMatch != true {
		respondError(w, http.StatusUnauthorized, "password does not match our records")
		return
	}

	var loginResponse struct {
		AccessToken string `json:"token"`
	}
	loginResponse.AccessToken = accessToken
	respondJSON(w, http.StatusAccepted, loginResponse)

}
