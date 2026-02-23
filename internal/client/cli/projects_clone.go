package cli

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func (c *Commands) ProjectsClone(args []string) error {
	if len(args) != 2 {
		return errors.New("projects clone takes 2 arguments: a project title and a destination directory path")
	}

	projectName := args[0]
	destPath, err := filepath.Abs(args[1])
	if err != nil {
		return fmt.Errorf("invalid destination path: %v", err)
	}

	projectRoot := filepath.Join(destPath, projectName)
	if err := os.MkdirAll(projectRoot, 0o755); err != nil {
		return fmt.Errorf("failed to create project root: %v", err)
	}

	assets, err := c.listAssetsForProject(projectName)
	if err != nil {
		return err
	}

	type downloadJob struct {
		Name string
		Tag  string
	}

	maxWorkers := 5
	jobs := make(chan downloadJob, len(assets))
	var wg sync.WaitGroup
	var firstErrMu sync.Mutex
	var firstErr error

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				dl, err := c.getDownloadUrl(projectName, job.Name)
				if err != nil {
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					firstErrMu.Unlock()
					continue
				}

				tagDir := projectRoot
				if job.Tag != "" {
					tagDir = filepath.Join(projectRoot, job.Tag)
				}
				os.MkdirAll(tagDir, 0o755)

				outPath := filepath.Join(tagDir, job.Name)
				fmt.Printf("Downloading %s...\n", job.Name)

				resp, err := http.Get(dl.UploadURL)
				if err != nil {
					continue
				}

				func() {
					defer resp.Body.Close()
					f, err := os.Create(outPath)
					if err != nil {
						return
					}
					defer f.Close()
					io.Copy(f, resp.Body)
				}()
			}
		}()
	}

	for name, tag := range assets {
		jobs <- downloadJob{Name: name, Tag: tag}
	}
	close(jobs)

	wg.Wait()
	http.DefaultTransport.(*http.Transport).CloseIdleConnections()

	return firstErr
}
