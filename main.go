package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bancey/ipmitool-api/internal/api"
	"github.com/bancey/ipmitool-api/internal/config"
	"github.com/bancey/ipmitool-api/internal/ipmi"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	executor := ipmi.NewCommandExecutor()
	server := api.NewServer(cfg, executor)

	log.Printf("Starting ipmitool-api on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
