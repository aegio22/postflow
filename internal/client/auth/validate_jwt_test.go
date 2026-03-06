package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	testUserID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	testSecret := "test-secret-key-12345"

	// Create a valid token for testing
	validToken, _ := MakeJWT(testUserID, testSecret)

	// Create an expired token
	expiredClaims := jwt.RegisteredClaims{
		Issuer:    "postflow",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC().Add(-2 * time.Hour)),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour)),
		Subject:   testUserID.String(),
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, _ := expiredTokenObj.SignedString([]byte(testSecret))

	tests := []struct {
		name           string
		tokenString    string
		tokenSecret    string
		isErr          bool
		expectedResult uuid.UUID
		expectedErr    string
	}{
		{
			name:           "success case basic",
			tokenString:    validToken,
			tokenSecret:    testSecret,
			isErr:          false,
			expectedResult: testUserID,
			expectedErr:    "",
		},
		{
			name:           "fail case wrong secret",
			tokenString:    validToken,
			tokenSecret:    "wrong-secret",
			isErr:          true,
			expectedResult: uuid.UUID{},
			expectedErr:    "token signature is invalid: signature is invalid",
		},
		{
			name:           "fail case malformed token",
			tokenString:    "not-a-valid-token",
			tokenSecret:    testSecret,
			isErr:          true,
			expectedResult: uuid.UUID{},
			expectedErr:    "token is malformed: token contains an invalid number of segments",
		},
		{
			name:           "fail case empty token",
			tokenString:    "",
			tokenSecret:    testSecret,
			isErr:          true,
			expectedResult: uuid.UUID{},
			expectedErr:    "token is malformed: token contains an invalid number of segments",
		},
		{
			name:           "fail case expired token",
			tokenString:    expiredToken,
			tokenSecret:    testSecret,
			isErr:          true,
			expectedResult: uuid.UUID{},
			expectedErr:    "token has invalid claims: token is expired",
		},
		{
			name:           "fail case invalid subject UUID",
			tokenString:    createTokenWithInvalidSubject(testSecret),
			tokenSecret:    testSecret,
			isErr:          true,
			expectedResult: uuid.Nil,
			expectedErr:    "invalid UUID length: 10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, errRaw := ValidateJWT(tt.tokenString, tt.tokenSecret)
			err := ""
			if errRaw != nil {
				err = errRaw.Error()
			}

			if result != tt.expectedResult || err != tt.expectedErr {
				t.Errorf("ValidateJWT test %v failed:\n\nExpected Result:%v\n Returned Result:%v\n Expected Error: %v\n Returned Error: %v\n",
					tt.name, tt.expectedResult, result, tt.expectedErr, err)
			}
		})
	}
}

// Helper function to create a token with an invalid subject for testing
func createTokenWithInvalidSubject(secret string) string {
	claims := jwt.RegisteredClaims{
		Issuer:    "postflow",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject:   "not-a-uuid",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, _ := token.SignedString([]byte(secret))
	return signedString
}
