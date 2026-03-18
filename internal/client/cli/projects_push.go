package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	activeWorkers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "postflow_active_workers",
		Help: "Number of goroutine workers currently processing S3 transfers",
	})
	transferDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "postflow_transfer_duration_seconds",
			Help:    "Duration of individual S3 upload and download operations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func init() {
	prometheus.MustRegister(activeWorkers)
	prometheus.MustRegister(transferDuration)
}

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
	err = c.CreateProject(projectArgs)
	if err != nil {
		return err
	}

	var assets []asset
	if err = filepath.WalkDir(sourcePath, helperParseLocalFiles(sourcePath, &assets)); err != nil {
		return fmt.Errorf("error walking project directory : %s", err)
	}

	maxWorkers := 5
	jobs := make(chan asset, len(assets))
	var wg sync.WaitGroup
	var firstErrMu sync.Mutex
	var firstErr error

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range jobs {
				activeWorkers.Inc()
				defer activeWorkers.Dec()

				fmt.Printf("Pushing %s...\n", a.AssetName)
				start := time.Now()
				err := c.UploadAsset([]string{projectName, a.Filepath, a.Tag})
				duration := time.Since(start)
				transferDuration.WithLabelValues("upload").Observe(duration.Seconds())

				if err != nil {
					slog.Error("S3 upload failed",
						"operation", "upload",
						"duration_ms", duration.Milliseconds(),
						"asset", a.AssetName,
						"error", err)
					firstErrMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					firstErrMu.Unlock()
					fmt.Printf("Upload failed for %s: %v\n", a.Filepath, err)
				} else {
					slog.Info("S3 upload completed",
						"operation", "upload",
						"duration_ms", duration.Milliseconds(),
						"asset", a.AssetName)
				}
			}
		}()
	}

	for _, a := range assets {
		jobs <- a
	}
	close(jobs)

	wg.Wait()

	if firstErr != nil {
		return firstErr
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
