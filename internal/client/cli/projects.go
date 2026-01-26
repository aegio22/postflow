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
	case "delete":
		return c.DeleteProject(args)
	case "delmem":
		return c.ProjectsDelmem(args)
	case "userlist":
		return c.ProjectsUserlist(args)
	default:
		return errors.New("unsupported projects command. run help for more info.")
	}

}
