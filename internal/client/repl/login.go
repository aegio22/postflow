package repl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) Login(args []string) error {
	if len(args) != 2 {
		return errors.New("sign up takes 2 arguments:email, password")
	}

	userCredentials := UserInfo{Email: args[0], Password: args[1]}
	requestBody, err := json.Marshal(userCredentials)
	if err != nil {
		return fmt.Errorf("error marshaling request body")
	}

	url := c.httpClient.BaseURL + routes.Login
	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		var errResp ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			// Fallback: couldn't parse error response
			return fmt.Errorf("login failed with status %d: %v", resp.StatusCode, err)
		}

		// Successfully parsed error response
		if errResp.Error != "" {
			return fmt.Errorf("login failed: %s", errResp.Error)
		}

		// Error response was empty
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}
	return nil

}
