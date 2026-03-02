package cli

import (
	"fmt"
	"os"

	"github.com/aegio22/postflow/internal/client/http"
	"github.com/aegio22/postflow/internal/server"
	"github.com/joho/godotenv"
)

// Command registry
type Commands struct {
	httpClient *http.HttpClient
}

type cliCommand struct {
	name        string
	description string
	arguments   []string
	callback    func([]string) error
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (c *Commands) getCommands() []cliCommand {
	//add commands here
	return []cliCommand{
		{
			name:        "register",
			description: "register a new user",
			arguments:   []string{"username", "email", "password"},
			callback:    c.SignUp,
		},
		{
			name:        "login",
			description: "login to an existing account",
			arguments:   []string{"email", "password"},
			callback:    c.Login,
		},
		{
			name:        "help",
			description: "get list of valid commands and their arguments",
			callback:    c.Help,
		},
		{
			name:        "serve",
			description: "initialize the postflow server. reads base server URL from environment variables",
			callback:    server.Run,
		},
		{
			name:     "projects",
			callback: c.Projects,
		},
		//not reachable from here. only in cmd registry for the help command's accuracy
		{
			name:        "projects clone",
			description: "clone a project from the database down to your local system",
			arguments:   []string{"project title", "destination filepath"},
			callback:    c.ProjectsClone,
		},
		{
			name:        "projects push",
			description: "push a project from your local system up to the database and cloud storage",
			arguments:   []string{"project title", "source filepath"},
			callback:    c.ProjectsPush,
		},

		{
			name:        "projects create",
			description: "create a new project",
			arguments:   []string{"project title", "'optional description'"},
			callback:    c.CreateProject,
		},
		{
			name:        "projects addmem",
			description: "add a new project member. Only available to admin or staff users within the project",
			arguments:   []string{"project title", "user email", "user status (admin, staff, or viewer)"},
			callback:    c.AddUserToProject,
		},
		{
			name:        "projects ls",
			description: "list all projects the logged in user is a member of",
			callback:    c.LsProjects,
		},
		{
			name:        "projects delete",
			description: "delete a project and all corresponding files from the database. Will only work if the logged in user is a project admin. **CANNOT BE UNDONE**",
			arguments:   []string{"project name"},
			callback:    c.DeleteProject,
		},
		{
			name:        "projects delmem",
			description: "delete a member from the given project",
			arguments:   []string{"project title", "user email"},
			callback:    c.DeleteAsset,
		},
		{
			name:        "projects userlist",
			description: "list all users and their statuses for a project you are a member of",
			arguments:   []string{"project title"},
			callback:    c.ProjectsUserlist,
		},
		//reachable
		{
			name:     "assets",
			callback: c.Assets,
		},
		//unreachable
		{
			name:        "assets upload",
			description: "upload new asset to project",
			arguments:   []string{"project name", "asset filepath", "optional tag"},
			callback:    c.UploadAsset,
		},
		{
			name:        "assets view",
			description: "download an asset",
			arguments:   []string{"project title", "asset filename"},
			callback:    c.ViewAsset,
		},
		{
			name:        "assets ls",
			description: "view all assets and their tags",
			arguments:   []string{"project title"},
			callback:    c.AssetsLs,
		},
		{
			name:        "assets delete",
			description: "delete an asset from cloud storage and database",
			arguments:   []string{"project title", "asset filename"},
			callback:    c.DeleteAsset,
		},
	}
}

func RunCLI() {
	_ = godotenv.Load()

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
	cmdMap := make(map[string]cliCommand)
	for _, cmd := range registry.getCommands() {
		cmdMap[cmd.name] = cmd
	}
	cmd, exists := cmdMap[cmdName]

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
