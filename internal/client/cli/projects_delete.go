package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) DeleteProject(args []string) error {
	if len(args) != 1 {
		return errors.New("projects delete takes one argument: project name")
	}
	projectName := args[0]
	url := c.httpClient.BaseURL + routes.Projects + "?project_name=" + url.QueryEscape(projectName)
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
	fmt.Printf("Project %s deleted successfully\n", projectName)

	return nil
}
