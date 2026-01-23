package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) ProjectsDelmem(args []string) error {
	if len(args) != 2 {
		return errors.New("projects delmem takes 2 arguments: a project title and user email")
	}

	projectName := args[0]
	userEmail := args[1]

	url := c.httpClient.BaseURL + routes.ProjectMembers + "?project_name=" + url.QueryEscape(projectName) + "&user_email=" + url.QueryEscape(userEmail)
	resp, err := c.httpClient.Delete(url)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("request failed: %s", errResp.Error)
	}
	fmt.Printf("User %s removed from project %s successfully\n", userEmail, projectName)

	return nil

}
