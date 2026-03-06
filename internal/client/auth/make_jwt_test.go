package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	testUserID := uuid.New()
	testSecret := "test-secret-key-12345"

	tests := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		isErr       bool
		expectedErr string
	}{
		{
			name:        "success case basic",
			userID:      testUserID,
			tokenSecret: testSecret,
			isErr:       false,
			expectedErr: "",
		},
		{
			name:        "success case with different UUID",
			userID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			tokenSecret: "another-secret",
			isErr:       false,
			expectedErr: "",
		},
		{
			name:        "success case with empty secret",
			userID:      testUserID,
			tokenSecret: "",
			isErr:       false,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errRaw := MakeJWT(tt.userID, tt.tokenSecret)
			err := ""
			if errRaw != nil {
				err = errRaw.Error()
			}

			if tt.isErr && err == "" {
				t.Errorf("MakeJWT test %v failed: expected error but got none", tt.name)
				return
			}

			if !tt.isErr && err != "" {
				t.Errorf("MakeJWT test %v failed: unexpected error: %v", tt.name, err)
				return
			}

			if err != tt.expectedErr {
				t.Errorf("MakeJWT test %v failed:\n Expected Error: %v\n Returned Error: %v\n",
					tt.name, tt.expectedErr, err)
				return
			}

			// If no error expected, validate the token structure
			if !tt.isErr && result != "" {
				// Parse the token to verify it's valid
				claims := jwt.RegisteredClaims{}
				token, parseErr := jwt.ParseWithClaims(result, &claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(tt.tokenSecret), nil
				})

				if parseErr != nil {
					t.Errorf("MakeJWT test %v failed: generated token is invalid: %v", tt.name, parseErr)
					return
				}

				if !token.Valid {
					t.Errorf("MakeJWT test %v failed: generated token is not valid", tt.name)
					return
				}

				// Verify claims
				if claims.Issuer != "postflow" {
					t.Errorf("MakeJWT test %v failed: expected issuer 'postflow', got '%v'", tt.name, claims.Issuer)
				}

				if claims.Subject != tt.userID.String() {
					t.Errorf("MakeJWT test %v failed: expected subject '%v', got '%v'", tt.name, tt.userID.String(), claims.Subject)
				}

				// Verify expiration is approximately 1 hour from now
				expectedExpiry := time.Now().UTC().Add(time.Hour)
				if claims.ExpiresAt.Time.Before(expectedExpiry.Add(-time.Minute)) || claims.ExpiresAt.Time.After(expectedExpiry.Add(time.Minute)) {
					t.Errorf("MakeJWT test %v failed: expiration time is not approximately 1 hour from now", tt.name)
				}
			}
		})
	}
}
