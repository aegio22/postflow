package cli

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) LsProjects(args []string) error {
	url := c.httpClient.BaseURL + routes.Projects
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

	var lsResp models.ProjectsLsResponse

	if err := json.NewDecoder(resp.Body).Decode(&lsResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Projects for user %s:\n\n", lsResp.UserName)
	for name, status := range lsResp.Projects {
		fmt.Printf("%s %s\n", name, status)
	}
	return nil

}
