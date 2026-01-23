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

func (c *Config) handlerLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo models.UserInfo
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

	refreshToken, err := c.DB.GetTokenFromUserID(ctx, user.ID)
	if err != nil {
		log.Printf("no refresh token found in DB for given user: %v", err)
		respondError(w, http.StatusUnauthorized, "no refresh token found in DB for given user")
		return
	}
	//Because of how sql.NullStrings work, RevokedAt is valid if the token HAS BEEN revoked
	if refreshToken.RevokedAt.Valid {
		log.Printf("refresh token is revoked for user %s", user.ID)
		respondError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		err = c.DB.RevokeToken(ctx, refreshToken.Token)
		if err != nil {
			log.Printf("error revoking expired refresh token: %v", err)
			respondError(w, http.StatusConflict, "error revoking expired refresh token")
			return
		}
		jwt, err := auth.MakeJWTSecret()
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusConflict, "old access token expired, and there is an error making a new one")
			return
		}
		_, err = c.DB.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{Token: jwt, UserID: user.ID, ExpiresAt: time.Now().AddDate(0, 0, 60)})
		if err != nil {
			log.Printf("error adding refresh token to database: %v", err)
			respondError(w, http.StatusBadRequest, "error adding refresh token to database")
			return
		}
		c.Env.JWT_SECRET = jwt
	}

	var loginResponse models.LoginResponse
	loginResponse.AccessToken = accessToken
	respondJSON(w, http.StatusAccepted, loginResponse)

}
