package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) ProjectsUserlist(args []string) error {
	if len(args) != 1 {
		return errors.New("projects userlist takes 1 argument: a project title")
	}

	projectName := args[0]

	url := c.httpClient.BaseURL + routes.ProjectMembers + "?project_name=" + url.QueryEscape(projectName)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("request failed: %s", errResp.Error)
	}
	var userlist models.ProjectsUserlistResponse

	if err := json.NewDecoder(resp.Body).Decode(&userlist); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Users in project %s:\n\n", projectName)
	for email, status := range userlist.Users {
		fmt.Printf("%s %s\n", email, status)
	}
	return nil
}
