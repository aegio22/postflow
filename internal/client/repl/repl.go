package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{}
}

func StartREPL() {

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Postflow > ")

		scanner.Scan()
		if scanner.Err() != nil {
			err := fmt.Errorf("error encountered: %v", scanner.Err())
			fmt.Println(err)
			continue
		}
		if scanner.Text() == "" {
			fmt.Println("no text in the scanner")
			continue
		}

		var args []string
		words := cleanInput(scanner.Text())
		commandName := words[0]

		if len(words) > 1 {
			args = words[1:]
		} else {
			args = []string{}
		}

		cmd, exists := getCommands()[commandName]

		if !exists {
			fmt.Println("Unknown command")
			continue

		} else {
			err := cmd.callback(args)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

	}

}
