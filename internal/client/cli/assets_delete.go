package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) DeleteAsset(args []string) error {
	if len(args) != 2 {
		return errors.New("assets delete takes 2 arguments: project title and asset filename")
	}
	projectName := args[0]
	assetName := args[1]
	url := c.httpClient.BaseURL + strings.Replace(routes.Assets, "{project_name}", projectName, 1) + "?project_name=" + url.QueryEscape(projectName) + "&asset_name=" + url.QueryEscape(assetName)
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
	fmt.Printf("Asset %s deleted from project %s successfully\n", assetName, projectName)

	return nil
}
