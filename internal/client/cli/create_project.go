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

func (c *Commands) CreateProject(args []string) error {
	if len(args) < 1 || len(args) > 2 {
		return errors.New("projects create takes 1-2 arguments: title, and an optional description.")
	}
	var description string
	var hasDesc bool
	if len(args) == 2 {
		description = args[1]
		hasDesc = true
	}
	project := models.ProjectRequest{
		Title:       args[0],
		Description: description,
	}
	requestBody, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("error marshaling request body")
	}
	url := c.httpClient.BaseURL + routes.Projects
	resp, err := c.httpClient.Post(url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		var errResp ErrorResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		return fmt.Errorf("project creation failed: %s", errResp.Error)
	}
	var projectResp models.ProjectResponse
	if err := json.NewDecoder(resp.Body).Decode(&projectResp); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Project created successfully!\n")
	fmt.Printf("Project ID: %s\n", projectResp.ID)
	fmt.Printf("Project Title: %s\n", projectResp.Title)
	if hasDesc == true {
		fmt.Printf("Project Description: %s\n", projectResp.Description)
	}
	fmt.Printf("Project Status: %s\n", projectResp.Status)
	fmt.Printf("Project Author: %s\n", projectResp.CreatedBy)

	return nil

}
