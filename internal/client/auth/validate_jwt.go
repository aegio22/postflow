package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	}
	token, err := jwt.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return uuid.UUID{}, err
	}
	subj, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(subj)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
