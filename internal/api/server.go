package api

import (
	"fmt"
	"net/http"

	"github.com/bancey/ipmitool-api/internal/config"
	"github.com/bancey/ipmitool-api/internal/ipmi"
)

func NewServer(cfg *config.Config, executor ipmi.Executor) *http.Server {
	h := &Handlers{Config: cfg, Executor: executor}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers", h.ListServers)
	mux.HandleFunc("GET /api/servers/{name}/power", h.GetPowerStatus)
	mux.HandleFunc("POST /api/servers/{name}/power", h.SetPowerState)
	mux.HandleFunc("GET /api/servers/{name}/sensors", h.GetSensors)
	mux.HandleFunc("GET /api/servers/{name}/chassis", h.GetChassisStatus)

	var handler http.Handler = mux
	handler = APIKeyAuth(cfg.API.APIKey, handler)

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.API.Port),
		Handler: handler,
	}
}
