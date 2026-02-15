package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/0xDevvvvv/Infra-Orchestrator/internal/models"
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/storage"
	"github.com/google/uuid"
)

type BuildHandler struct {
	store *storage.BuildStore
}

type CreateBuildRequest struct {
	RepoURL string `json:"repo_url"`
	Branch  string `json:"branch"`
}

type CreateBuildResponse struct {
	ID     string             `json:"id"`
	Status models.BuildStatus `json:"status"`
}

func NewBuildHandler(store *storage.BuildStore) *BuildHandler {
	return &BuildHandler{
		store: store,
	}
}

func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateBuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.RepoURL == "" || req.Branch == "" {
		http.Error(w, "Repo URL and Branch are required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()

	build := &models.Build{
		ID:        id,
		RepoURL:   req.RepoURL,
		Branch:    req.Branch,
		Status:    models.Pending,
		CreatedAt: time.Now(),
	}

	h.store.Save(build)

	response := CreateBuildResponse{
		ID:     id,
		Status: build.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
