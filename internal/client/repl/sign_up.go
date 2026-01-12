package repl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/client/auth"
)

type UserInfo struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
}

func (c *Commands) SignUp(args []string) error {
	if len(args) != 3 {
		return errors.New("sign up takes 3 arguments: username, email, password")
	}

	hashedPassword, err := auth.HashPassword(args[2])
	if err != nil {
		return err
	}
	newUser := UserInfo{Username: args[0], Email: args[1], HashedPassword: hashedPassword}
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		return fmt.Errorf("error marshaling request body")
	}

	url := c.httpClient.BaseURL + "/api/signup"
	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			// Fallback: couldn't parse error response
			return fmt.Errorf("signup failed with status %d: %w", resp.StatusCode, err)
		}

		// Successfully parsed error response
		if errResp.Error != "" {
			return fmt.Errorf("signup failed: %s", errResp.Error)
		}

		// Error response was empty
		return fmt.Errorf("signup failed with status %d", resp.StatusCode)
	}
	return nil

}
