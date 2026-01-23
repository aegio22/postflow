package cli

import "errors"

func (c *Commands) Assets(args []string) error {
	if len(args) == 0 {
		return errors.New("no subcommand provided for assets")
	}
	sub := args[0]
	args = args[1:]
	switch sub {
	case "upload":
		return c.UploadAsset(args)
	case "view":
		return c.ViewAsset(args)
	case "ls":
		return c.AssetsLs(args)
	default:
		return errors.New("unsupported assets command. run assets help for more info.")
	}

	return nil
}
