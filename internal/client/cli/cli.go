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
	Error   string `json:"error"`             // Main error message
	Message string `json:"message,omitempty"` // Optional detailed message
	Code    string `json:"code,omitempty"`    // Optional error code (e.g., "INVALID_EMAIL")
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
