package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	crypt := make([]byte, 32)
	bytesWritten, err := rand.Read(crypt)
	if err != nil {
		return "", fmt.Errorf("error generating refresh token: %v", err)
	}
	if bytesWritten != 32 {
		return "", errors.New("insufficient data size written")
	}

	token := hex.EncodeToString(crypt)
	return token, nil

}
