package repl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/routes"
)

type UserInfo struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"hashed_password"`
}

type SignUpResponse struct {
	UserID  string `json:"id"`
	Token   string `json:"access_token"`
	Message string `json:"message"`
}

func (c *Commands) SignUp(args []string) error {
	if len(args) != 3 {
		return errors.New("sign up takes 3 arguments: username, email, password")
	}

	newUser := UserInfo{Username: args[0], Email: args[1], Password: args[2]}
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		return fmt.Errorf("error marshaling request body")
	}

	url := c.httpClient.BaseURL + routes.SignUp
	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("signup failed: %s", errResp.Error)
	}

	var signupResp SignUpResponse
	if err := json.NewDecoder(resp.Body).Decode(&signupResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	c.httpClient.SetAuthToken(signupResp.Token)
	fmt.Printf("Account created successfully!\n")
	fmt.Printf("User ID: %s\n", signupResp.UserID)
	fmt.Printf("Logged in automatically\n")
	return nil

}
