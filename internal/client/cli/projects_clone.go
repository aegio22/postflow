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
		return errors.New("projects push takes 2 arguments: a project title and a destination directory path")
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
	jobs := make(chan downloadJob)
	var wg sync.WaitGroup
	var firstErrMu sync.Mutex
	var firstErr error

	// start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				// Get presigned download URL
				dl, err := c.getDownloadUrl(projectName, job.Name)
				if err != nil {
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					firstErrMu.Unlock()
					fmt.Printf("failed to get download URL for %s: %v\n", job.Name, err)
					continue
				}

				// Compute tag directory
				tagDir := projectRoot
				if job.Tag != "" {
					tagDir = filepath.Join(projectRoot, job.Tag)
				}
				if err := os.MkdirAll(tagDir, 0o755); err != nil {
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = fmt.Errorf("error making subdirectory %s: %s", tagDir, err)
					}
					firstErrMu.Unlock()
					fmt.Printf("failed to create dir %s: %v\n", tagDir, err)
					continue
				}

				outPath := filepath.Join(tagDir, job.Name)
				fmt.Printf("Downloading file %s\n", job.Name)

				// Download from S3
				resp, err := http.Get(dl.UploadURL)
				if err != nil {
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = fmt.Errorf("error downloading asset %s: %s", job.Name, err)
					}
					firstErrMu.Unlock()
					fmt.Printf("download failed for %s: %v\n", job.Name, err)
					continue
				}
				func() {
					defer resp.Body.Close()
					f, err := os.Create(outPath)
					if err != nil {
						firstErrMu.Lock()
						if firstErr == nil {
							firstErr = fmt.Errorf("error creating file at %s: %s", outPath, err)
						}
						firstErrMu.Unlock()
						fmt.Printf("create file failed for %s: %v\n", outPath, err)
						return
					}
					defer f.Close()

					if _, err := io.Copy(f, resp.Body); err != nil {
						firstErrMu.Lock()
						if firstErr == nil {
							firstErr = fmt.Errorf("could not write file %s: %s", outPath, err)
						}
						firstErrMu.Unlock()
						fmt.Printf("write failed for %s: %v\n", outPath, err)
						return
					}
					fmt.Println("Download successful!")
				}()
			}
		}()
	}

	// enqueue jobs
	for name, tag := range assets {
		jobs <- downloadJob{Name: name, Tag: tag}
	}
	close(jobs)

	wg.Wait()

	if firstErr != nil {
		return firstErr
	}

	return nil
}
