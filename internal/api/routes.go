package api

import (
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/api/handlers"
)

func (s *Server) registerRoutes() {

	buildHandler := handlers.NewBuildHandler(s.store, s.queue)
	
	s.router.HandleFunc("/builds", buildHandler.CreateBuild)
	s.router.HandleFunc("/builds/", buildHandler.GetBuild)
}
