package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/0xDevvvvv/Infra-Orchestrator/internal/build"
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/models"
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/queue"
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/storage"
)

type Server struct {
	router *http.ServeMux
	store  *storage.BuildStore
	queue  *queue.BuildQueue
	runner *build.Runner
}

func NewServer() *Server {
	s := &Server{
		router: http.NewServeMux(),
		store:  storage.NewBuildStore(),
		queue:  queue.NewBuildQueue(100),
		runner: build.NewRunner("tmp", "artifacts", 5*time.Minute),
	}

	s.registerRoutes()
	go s.StartWorker()
	return s
}

func (s *Server) StartWorker() {
	for {
		id := s.queue.Dequeue()

		build, ok := s.store.Get(id)
		if !ok {
			continue
		}

		build.Status = models.Running

		err := s.runner.Run(build)
		if err != nil {
			fmt.Println("RUNNER ERROR:", err)
			build.Status = models.Failed
			continue
		}
		build.Status = models.Success
	}
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
