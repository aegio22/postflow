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

func (c *Commands) AddUserToProject(args []string) error {
	if len(args) != 3 {
		return errors.New("projects add takes 3 arguments: project name, user email, and user status")
	}

	//parse user status arg for validity check
	if args[2] != "admin" && args[2] != "staff" && args[2] != "viewer" {
		return errors.New("invalid user status provided. Users can either be 'admin', 'staff', or 'viewer'")
	}
	userAdd := models.AddUserRequest{
		ProjectName: args[0],
		UserEmail:   args[1],
		UserStatus:  args[2],
	}
	requestBody, err := json.Marshal(userAdd)
	if err != nil {
		return errors.New("error marshaling request body")
	}
	url := c.httpClient.BaseURL + routes.ProjectMembers
	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("project member addition failed: %s", errResp.Error)
	}

	var responseBody models.ProjectMemberAddResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	fmt.Printf("Member added succesfully to %s with %s permissions", responseBody.ProjectName, responseBody.UserStatus)
	return nil
}
