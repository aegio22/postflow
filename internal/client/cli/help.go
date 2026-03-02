package cli

import (
	"fmt"
	"strings"
)

func (c *Commands) Help(args []string) error {
	fmt.Println("Welcome to Postflow!")
	fmt.Println("This is the base CLI tool for the project management platform.")
	fmt.Println("For more info on integrating your own UI with Postflow, view the documentation on github.")
	fmt.Print("Usage: postflow [command]\n\n")

	cmds := c.getCommands()

	lastPrefix := ""
	for _, cmd := range cmds {
		currentPrefix := strings.SplitN(cmd.name, " ", 2)[0]
		if currentPrefix != lastPrefix {
			if lastPrefix != "" {
				fmt.Println()
			}
			lastPrefix = currentPrefix
		}

		fmt.Printf("%v:\n", cmd.name)
		var sb strings.Builder
		if len(cmd.arguments) > 0 {
			sb.WriteString("Arguments: ")
			for i := 0; i < len(cmd.arguments); i++ {
				arg := fmt.Sprintf("<%v>", cmd.arguments[i])
				sb.WriteString(arg)
			}
		}
		argStr := sb.String()
		descString := cmd.description
		if argStr != "" {
			descString += ("\n" + "	" + argStr)
		}
		fmt.Println("	" + descString)
	}
	return nil
}
