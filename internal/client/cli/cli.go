package cli

import (
	"fmt"
	"os"

	"github.com/aegio22/postflow/internal/client/http"
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
	}
}

func RunCLI() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "No arguments or commands provided")
		os.Exit(2)
	}
	cmdName := os.Args[1]
	args := os.Args[2:]

	client := http.CreateHttpClient("")
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
