package cli

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
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
				activeWorkers.Inc()
				defer activeWorkers.Dec()

				start := time.Now()
				dl, err := c.getDownloadUrl(projectName, job.Name)
				if err != nil {
					duration := time.Since(start)
					slog.Error("failed to get download URL",
						"operation", "download",
						"duration_ms", duration.Milliseconds(),
						"asset", job.Name,
						"error", err)
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					firstErrMu.Unlock()
					transferDuration.WithLabelValues("download").Observe(duration.Seconds())
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
				duration := time.Since(start)
				if err != nil {
					slog.Error("S3 download failed",
						"operation", "download",
						"duration_ms", duration.Milliseconds(),
						"asset", job.Name,
						"error", err)
					transferDuration.WithLabelValues("download").Observe(duration.Seconds())
					continue
				}

				downloadErr := func() error {
					defer resp.Body.Close()
					f, err := os.Create(outPath)
					if err != nil {
						return err
					}
					defer f.Close()
					_, err = io.Copy(f, resp.Body)
					return err
				}()

				transferDuration.WithLabelValues("download").Observe(time.Since(start).Seconds())

				if downloadErr != nil {
					slog.Error("failed to save downloaded file",
						"operation", "download",
						"duration_ms", duration.Milliseconds(),
						"asset", job.Name,
						"error", downloadErr)
				} else {
					slog.Info("S3 download completed",
						"operation", "download",
						"duration_ms", duration.Milliseconds(),
						"asset", job.Name)
				}
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
