package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aegio22/postflow/internal/client/repl"
	"github.com/aegio22/postflow/internal/database"
	"github.com/google/uuid"
)

type DBUserResponse struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
}

func (c *Config) handlerSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var userInfo repl.UserInfo
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&userInfo)
	if err != nil {
		log.Printf("error decoding request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUser, err := c.DB.CreateUser(ctx, database.CreateUserParams{
		Username: userInfo.Username, Email: userInfo.Email, HashedPassword: userInfo.HashedPassword,
	})
	if err != nil {
		log.Printf("error registering user: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respUser := DBUserResponse{
		ID:        newUser.ID,
		Username:  newUser.Username,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}

	responseBody, err := json.Marshal(respUser)
	if err != nil {
		log.Printf("error marshaling response body: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(responseBody)

}
