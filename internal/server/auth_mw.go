package server

import (
	"context"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "user_id"

// requireAuth validates the JWT token and adds the user ID to the request context
func (c *Config) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "missing or invalid authorization header")
			return
		}

		userID, err := auth.ValidateJWT(token, c.Env.JWT_SECRET)
		if err != nil {
			respondError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

// getUserID extracts the user ID from the request context
func getUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
