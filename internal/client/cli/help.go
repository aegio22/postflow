package cli

import "fmt"

func (c *Commands) Help(args []string) error {
	fmt.Println("Welcome to Postflow!")
	fmt.Print("Usage:\n\n")

	for _, cmd := range c.getCommands() {
		line := fmt.Sprintf("%v: %v", cmd.name, cmd.description)
		fmt.Println(line)
	}
	return nil

}
