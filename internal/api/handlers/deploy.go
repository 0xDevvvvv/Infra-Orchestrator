package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type DeployHandler struct {
	artifactDir string
}

func NewDeployHandler(artifactDir string) *DeployHandler {
	return &DeployHandler{
		artifactDir: artifactDir,
	}
}

func (h *DeployHandler) ServeDeployment(w http.ResponseWriter, r *http.Request) {

	host := r.Host
	host = strings.Split(host, ":")[0]

	parts := strings.Split(host, ".")

	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}

	buildID := parts[0]

	dir := filepath.Join(h.artifactDir, buildID)

	// Check if build exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	fs := http.FileServer(http.Dir(dir))
	fs.ServeHTTP(w, r)
}
