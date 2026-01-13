package auth

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	expireTime := time.Hour
	claims := jwt.RegisteredClaims{
		Issuer:    "postflow",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expireTime)),
		Subject:   userID.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("error signing string: %v", err)
		return "", err
	}
	return signedString, nil
}
