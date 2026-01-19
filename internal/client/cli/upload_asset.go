package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aegio22/postflow/internal/client/models"
	"github.com/aegio22/postflow/internal/routes"
)

func (c *Commands) UploadAsset(args []string) error {
	if len(args) != 4 {
		return errors.New("assets upload takes 4 arguments: ProjectName, AssetPath, and Tag")
	}
	
	projectName := args[0]
	assetPath := args[2]
	tag := args[3]
	assetName := filepath.Base(assetPath)
	file, err := os.Open(assetPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %w", err)
	}

	fileName := filepath.Base(assetPath)
	fileSize := fileInfo.Size()
	url := c.httpClient.BaseURL + strings.Replace(routes.Assets, "{project_name}", projectName, 1)

	assetRequest := models.AssetRequest{
		ProjectName: projectName,
		AssetName:   assetName,
		Filepath:    assetPath,
		Tag:         tag,
	}

	requestBody, err := json.Marshal(assetRequest)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %s", err)
	}

	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("asset upload failed: %s", errResp.Error)
	}

	var assetResp models.UploadAssetResponse
	err = json.NewDecoder(resp.Body).Decode(&assetResp)
	if err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}

	fmt.Printf("Asset created: %s\n", assetResp.AssetID)
	fmt.Printf("Uploading %s (%.2f MB)...\n", fileName, float64(fileSize)/(1024*1024))

	file.Seek(0, 0) // reset file pointer
	putReq, err := http.NewRequest(http.MethodPut, assetResp.UploadURL, file)
	if err != nil {
		return fmt.Errorf("failed to create upload request: %w", err)
	}

	putReq.ContentLength = fileSize
	putReq.Header.Set("Content-Type", "application/octet-stream")
	client := &http.Client{}
	putResp, err := client.Do(putReq)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusOK && putResp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(putResp.Body)
		return fmt.Errorf("S3 upload failed (%d): %s", putResp.StatusCode, string(body))
	}

	fmt.Printf("Upload complete!\n")
	fmt.Printf("Asset ID: %s\n", assetResp.AssetID)
	fmt.Printf("S3 Key: %s\n", assetResp.S3Key)

	return nil
}
