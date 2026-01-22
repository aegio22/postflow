package cli

import (
	"errors"
)

func (c *Commands) Projects(args []string) error {
	if len(args) == 0 {
		return errors.New("no subcommand provided for projects")
	}
	sub := args[0]
	args = args[1:]
	switch sub {
	case "create":
		return c.CreateProject(args)
	case "addmem":
		return c.AddUserToProject(args)
	case "ls":
		return c.LsProjects(args)
	default:
		return errors.New("unsupported projects command. run projects help for more info.")
	}

}
