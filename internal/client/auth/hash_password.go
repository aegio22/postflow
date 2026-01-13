package auth

import (
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		errMessage := fmt.Sprintf("error hashing password: %v", err)
		return "", errors.New(errMessage)
	}
	return hashedPassword, nil
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		errMessage := fmt.Sprintf("error comparing password with hash: %v", err)
		return false, errors.New(errMessage)
	}
	return match, nil
}

