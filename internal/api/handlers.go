package api

import (
	"encoding/json"
	"net/http"

	"github.com/bancey/ipmitool-api/internal/config"
	"github.com/bancey/ipmitool-api/internal/ipmi"
)

type Handlers struct {
	Config   *config.Config
	Executor ipmi.Executor
}

type serverInfo struct {
	Name string `json:"name"`
	Host string `json:"host"`
}

func (h *Handlers) ListServers(w http.ResponseWriter, r *http.Request) {
	servers := make([]serverInfo, len(h.Config.Servers))
	for i, s := range h.Config.Servers {
		servers[i] = serverInfo{Name: s.Name, Host: s.Host}
	}
	writeJSON(w, http.StatusOK, servers)
}

type powerRequest struct {
	Action string `json:"action"`
}

func (h *Handlers) GetPowerStatus(w http.ResponseWriter, r *http.Request) {
	server, ok := h.findServer(w, r)
	if !ok {
		return
	}

	status, err := ipmi.GetPowerStatus(r.Context(), h.Executor, server)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (h *Handlers) SetPowerState(w http.ResponseWriter, r *http.Request) {
	server, ok := h.findServer(w, r)
	if !ok {
		return
	}

	var req powerRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1024)).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := ipmi.SetPowerState(r.Context(), h.Executor, server, req.Action); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"result": "ok"})
}

func (h *Handlers) GetSensors(w http.ResponseWriter, r *http.Request) {
	server, ok := h.findServer(w, r)
	if !ok {
		return
	}

	sensors, err := ipmi.GetSensors(r.Context(), h.Executor, server)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, sensors)
}

func (h *Handlers) GetChassisStatus(w http.ResponseWriter, r *http.Request) {
	server, ok := h.findServer(w, r)
	if !ok {
		return
	}

	status, err := ipmi.GetChassisStatus(r.Context(), h.Executor, server)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (h *Handlers) findServer(w http.ResponseWriter, r *http.Request) (*config.Server, bool) {
	name := r.PathValue("name")
	server, err := h.Config.FindServer(name)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return nil, false
	}
	return server, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
