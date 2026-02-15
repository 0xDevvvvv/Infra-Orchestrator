package main

import (
	"github.com/0xDevvvvv/Infra-Orchestrator/internal/api"
)

func main() {
	s := api.NewServer()
	println("Server Starting on port : 8080")
	s.Start(":8080")
}
