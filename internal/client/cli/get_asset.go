package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aegio22/postflow/internal/client/models"
)

func (c *Commands) GetAsset(args []string) error {
	if len(args) != 2 {
		return errors.New("assets get takes 2 arguments: project_name and asset_name (eg. assets get project1 final_cut_v7.mov)")
	}
	projectName := args[0]
	assetName := args[1]

	url := fmt.Sprintf("%s/assets/view?project_name=%s&asset_name=%s",
		c.httpClient.BaseURL, projectName, assetName)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("asset upload failed: %s", errResp.Error)
	}

	var responseBody models.AssetResponse
	downloadUrl := responseBody.UploadURL
	expiresIn := responseBody.ExpiresIn / 60
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	fmt.Println("Asset fetched succesfully!")
	fmt.Printf("URL (Expires in %v minutes):\n", expiresIn)
	fmt.Println(downloadUrl)

	return nil
}
