package cli

import (
	"fmt"
	"os"

	"github.com/aegio22/postflow/internal/client/http"
	"github.com/aegio22/postflow/internal/server"
)

// Command registry
type Commands struct {
	httpClient *http.HttpClient
}

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (c *Commands) getCommands() map[string]cliCommand {
	//add commands here
	return map[string]cliCommand{
		"register": {
			name:        "register",
			description: "register a new user",
			callback:    c.SignUp,
		},
		"login": {
			name:        "login",
			description: "login with email and password",
			callback:    c.Login,
		},
		"serve": {
			name:        "serve",
			description: "initialize the postflow server",
			callback:    server.Run,
		},
		"projects": {
			name:        "projects",
			description: "followed by projects subcommands",
			callback:    c.Projects,
		},
		//not reachable from here. only in cmd registry for the help command's accuracy
		"projects create": {
			name:        "projects create",
			description: "create a new project",
			callback:    c.CreateProject,
		},
		"projects addmem": {
			name:        "projects addmem",
			description: "add a new project member by project title, user email, and user status (admin, staff, or viewer)",
			callback:    c.AddUserToProject,
		},
		"projects ls": {
			name:        "projects ls",
			description: "list all projects the logged in user is a member of",
			callback:    c.LsProjects,
		},
		"projects delete": {
			name:        "projects delete",
			description: "delete a project from the database by project title. Will only work if the logged in user is a project admin",
			callback:    c.DeleteProject,
		},
		//reachable
		"assets": {
			name:        "assets",
			description: "followed by assets subcommands",
			callback:    c.Assets,
		},
		//unreachable
		"assets upload": {
			name:        "assets upload",
			description: "upload new asset to project by project title",
			callback:    c.UploadAsset,
		},
		"assets view": {
			name:        "assets view",
			description: "download an asset by project title and asset filename",
			callback:    c.ViewAsset,
		},
		"assets ls": {
			name:        "assets ls",
			description: "view all assets and their tags for a given project title",
			callback:    c.AssetsLs,
		},
		"assets delete": {
			name:        "assets delete",
			description: "delete an asset from storage and db via a given project title and asset filename",
			callback:    c.DeleteAsset,
		},
	}
}

func RunCLI() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "No arguments or commands provided")
		os.Exit(2)
	}
	cmdName := os.Args[1]
	args := os.Args[2:]

	client := http.CreateHttpClient()
	registry := Commands{
		httpClient: client,
	}
	cmd, exists := registry.getCommands()[cmdName]

	if !exists {
		fmt.Fprint(os.Stderr, "Unknown command")
		os.Exit(2)
	} else {
		err := cmd.callback(args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Callback error: %v", err)
		}
	}

}
