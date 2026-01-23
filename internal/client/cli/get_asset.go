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

func (c *Commands) ViewAsset(args []string) error {
	if len(args) != 2 {
		return errors.New("assets get takes 2 arguments: project_name and asset_name (eg. assets get project1 final_cut_v7.mov)")
	}
	projectName := args[0]
	assetName := args[1]

	url := c.httpClient.BaseURL + routes.ViewAssets + "?project_name=" + url.QueryEscape(projectName) + "&asset_name=" + url.QueryEscape(assetName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("asset view failed: %s", errResp.Error)
	}

	var responseBody models.AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	// NOW compute from the decoded struct
	downloadURL := responseBody.UploadURL
	expiresMinutes := responseBody.ExpiresIn / 60

	fmt.Println("Asset fetched succesfully!")
	fmt.Printf("URL (Expires in %v minutes):\n", expiresMinutes)
	fmt.Println(downloadURL)

	return nil
}
