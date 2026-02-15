package api

import (
	"net/http"

	"github.com/0xDevvvvv/Infra-Orchestrator/internal/storage"
)

type Server struct {
	router *http.ServeMux
	store  *storage.BuildStore
}

func NewServer() *Server {
	s := &Server{
		router: http.NewServeMux(),
		store:  storage.NewBuildStore(),
	}

	s.registerRoutes()

	return s
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
