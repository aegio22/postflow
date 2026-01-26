package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type asset struct {
	AssetName string
	Tag       string
	Filepath  string
}

func (c *Commands) ProjectsPush(args []string) error {
	if len(args) != 2 {
		return errors.New("projects push takes 2 arguments: a new project title and a source directory path")
	}

	projectName := args[0]
	sourcePath, err := filepath.Abs(args[1])
	if err != nil {
		return fmt.Errorf("invalid source path: %v", err)
	}

	var projectArgs []string
	projectArgs = append(projectArgs, projectName)
	//project creation
	err = c.CreateProject(projectArgs)
	if err != nil {
		return err
	}
	//asset uploads, local parsing
	var assets []asset
	if err = filepath.WalkDir(sourcePath, helperParseLocalFiles(sourcePath, &assets)); err != nil {
		return fmt.Errorf("error walking project directory : %s", err)
	}
	for _, a := range assets {
		if err := c.UploadAsset([]string{projectName, a.Filepath, a.Tag}); err != nil {
			return err
		}
	}

	return nil

}

func helperParseLocalFiles(root string, assets *[]asset) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		dir := filepath.Dir(rel)
		tag := ""
		if dir != "." {
			parts := strings.Split(dir, string(os.PathSeparator))
			tag = parts[len(parts)-1] // just the top-level folder name as tag
		}

		*assets = append(*assets, asset{
			AssetName: d.Name(),
			Tag:       tag,
			Filepath:  path,
		})

		return nil
	}
}
