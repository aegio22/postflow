package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) Login(args []string) error {
	if len(args) != 2 {
		return errors.New("sign up takes 2 arguments:email, password")
	}

	userCredentials := models.UserInfo{Email: args[0], Password: args[1]}
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
	var loginResponse struct {
		AccessToken string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}
	c.httpClient.SetSession(loginResponse.AccessToken)
	return nil

}
