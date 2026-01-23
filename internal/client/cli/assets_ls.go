package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) AssetsLs(args []string) error {
	if len(args) != 1 {
		return errors.New("assets ls takes one argument: project name")
	}

	projectName := args[0]
	query := fmt.Sprintf("?project_name=%s", projectName)
	url := c.httpClient.BaseURL + strings.Replace(routes.Assets, "{project_name}", projectName, 1) + query

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("error performing request: %v", err)
	}
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("request failed: %s", errResp.Error)
	}
	var lsResp models.AssetsLsResponse

	if err := json.NewDecoder(resp.Body).Decode(&lsResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Assets for project '%s':\n", projectName)
	for name, tag := range lsResp.Assets {
		fmt.Printf("%s %s\n", name, tag)
	}
	return nil
}
