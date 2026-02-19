package api

import (
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/api/handlers"
)

func (s *Server) registerRoutes() {

	buildHandler := handlers.NewBuildHandler(s.store, s.queue)
	deployHandler := handlers.NewDeployHandler("artifacts")

	s.router.HandleFunc("/builds", buildHandler.CreateBuild)
	s.router.HandleFunc("/builds/", buildHandler.GetBuild)

	s.router.HandleFunc("/", deployHandler.ServeDeployment)

}
