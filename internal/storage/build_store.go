package storage

import (
	"sync"

	"github.com/0xDevvvvv/Infra-Orchestrator/internal/models"
)

type BuildStore struct {
	builds map[string]*models.Build
	mu     sync.RWMutex
}

func NewBuildStore() *BuildStore {
	return &BuildStore{
		builds: make(map[string]*models.Build),
	}
}

func (bs *BuildStore) Save(build *models.Build) {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	bs.builds[build.ID] = build
}
