package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("header not found: authorization")
	}
	if !strings.HasPrefix(header, "Bearer ") {
		return "", errors.New("bearer token not found")
	}
	return strings.TrimPrefix(header, "Bearer "), nil
}
