package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/0xDevvvvv/Infra-Orchestrator/internal/models"
)

type Runner struct {
	tmpDir      string
	artifactDir string
	timeout     time.Duration
}

func NewRunner(tmpDir string, artifactDir string, timeout time.Duration) *Runner {
	return &Runner{
		tmpDir:      tmpDir,
		artifactDir: artifactDir,
		timeout:     timeout,
	}
}

func (r *Runner) Run(build *models.Build) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	workspace := filepath.Join(r.tmpDir, build.ID)

	if err := os.MkdirAll(workspace, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(workspace)

	if err := r.runCommand(ctx, workspace, "git", "clone", build.RepoURL, "."); err != nil {
		return err
	}

	if err := r.runCommand(ctx, workspace, "git", "checkout", build.Branch); err != nil {
		return err
	}

	if err := r.runCommand(ctx, workspace, "npm", "install"); err != nil {
		return err
	}

	if err := r.runCommand(ctx, workspace, "npm", "run", "build"); err != nil {
		return err
	}

	return r.saveArtifacts(workspace, build)

}

func (r *Runner) runCommand(ctx context.Context, dir string, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (r *Runner) saveArtifacts(workspace string, build *models.Build) error {

	possibleDirs := []string{"build", "dist"}
	var outputDir string

	for _, dir := range possibleDirs {
		path := filepath.Join(workspace, dir)
		if _, err := os.Stat(path); err == nil {
			outputDir = path
			break
		}
	}

	if outputDir == "" {
		return fmt.Errorf("no build output directory found")
	}

	target := filepath.Join(r.artifactDir, build.ID)

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	return copyDir(outputDir, target)

}

func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = io.Copy(dstFile, srcFile)
		return err
	})
}
