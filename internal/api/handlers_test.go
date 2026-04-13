package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bancey/ipmitool-api/internal/config"
)

type apiMockExecutor struct {
	output string
	err    error
}

func (m *apiMockExecutor) Execute(_ context.Context, _ *config.Server, args ...string) (string, error) {
	return m.output, m.err
}

var testConfig = &config.Config{
	API: config.APIConfig{Port: 8080, APIKey: ""},
	Servers: []config.Server{
		{Name: "server1", Host: "10.0.0.1", Port: 623, Username: "admin", Password: "pass"},
		{Name: "server2", Host: "10.0.0.2", Port: 623, Username: "admin", Password: "pass"},
	},
}

func newTestHandlers(output string, err error) *Handlers {
	return &Handlers{
		Config:   testConfig,
		Executor: &apiMockExecutor{output: output, err: err},
	}
}

func TestListServers(t *testing.T) {
	h := newTestHandlers("", nil)
	req := httptest.NewRequest("GET", "/api/servers", nil)
	rr := httptest.NewRecorder()

	h.ListServers(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var servers []serverInfo
	if err := json.NewDecoder(rr.Body).Decode(&servers); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(servers))
	}
	if servers[0].Name != "server1" {
		t.Errorf("expected server1, got %q", servers[0].Name)
	}
}

func TestGetPowerStatus(t *testing.T) {
	h := newTestHandlers("Chassis Power is on", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/power", h.GetPowerStatus)

	req := httptest.NewRequest("GET", "/api/servers/server1/power", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var result map[string]string
	json.NewDecoder(rr.Body).Decode(&result)
	if result["status"] != "on" {
		t.Errorf("expected status 'on', got %q", result["status"])
	}
}

func TestGetPowerStatusNotFound(t *testing.T) {
	h := newTestHandlers("", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/power", h.GetPowerStatus)

	req := httptest.NewRequest("GET", "/api/servers/nonexistent/power", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestGetPowerStatusError(t *testing.T) {
	h := newTestHandlers("", fmt.Errorf("connection refused"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/power", h.GetPowerStatus)

	req := httptest.NewRequest("GET", "/api/servers/server1/power", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestSetPowerState(t *testing.T) {
	h := newTestHandlers("", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/servers/{name}/power", h.SetPowerState)

	body := strings.NewReader(`{"action":"on"}`)
	req := httptest.NewRequest("POST", "/api/servers/server1/power", body)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestSetPowerStateInvalidBody(t *testing.T) {
	h := newTestHandlers("", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/servers/{name}/power", h.SetPowerState)

	body := strings.NewReader(`not json`)
	req := httptest.NewRequest("POST", "/api/servers/server1/power", body)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestSetPowerStateInvalidAction(t *testing.T) {
	h := newTestHandlers("", nil)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/servers/{name}/power", h.SetPowerState)

	body := strings.NewReader(`{"action":"destroy"}`)
	req := httptest.NewRequest("POST", "/api/servers/server1/power", body)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestGetSensors(t *testing.T) {
	output := `CPU Temp         | 45.000     | degrees C  | ok
Fan1             | 3200.000   | RPM        | ok`

	h := newTestHandlers(output, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/sensors", h.GetSensors)

	req := httptest.NewRequest("GET", "/api/servers/server1/sensors", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var sensors []map[string]string
	json.NewDecoder(rr.Body).Decode(&sensors)
	if len(sensors) != 2 {
		t.Errorf("expected 2 sensors, got %d", len(sensors))
	}
}

func TestGetChassisStatus(t *testing.T) {
	output := `System Power         : on
Power Overload       : false
Drive Fault          : false
Cooling/Fan Fault    : false`

	h := newTestHandlers(output, nil)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/chassis", h.GetChassisStatus)

	req := httptest.NewRequest("GET", "/api/servers/server1/chassis", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	var result map[string]any
	json.NewDecoder(rr.Body).Decode(&result)
	if result["power_on"] != true {
		t.Errorf("expected power_on true, got %v", result["power_on"])
	}
}

func TestGetChassisStatusError(t *testing.T) {
	h := newTestHandlers("", fmt.Errorf("timeout"))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/servers/{name}/chassis", h.GetChassisStatus)

	req := httptest.NewRequest("GET", "/api/servers/server1/chassis", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestFullServerWithAuth(t *testing.T) {
	cfg := &config.Config{
		API: config.APIConfig{Port: 8080, APIKey: "secret"},
		Servers: []config.Server{
			{Name: "srv1", Host: "10.0.0.1", Port: 623, Username: "admin", Password: "pass"},
		},
	}
	mock := &apiMockExecutor{output: "Chassis Power is on"}
	srv := NewServer(cfg, mock)

	// Request without API key
	req := httptest.NewRequest("GET", "/api/servers", nil)
	rr := httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 without API key, got %d", rr.Code)
	}

	// Request with valid API key
	req = httptest.NewRequest("GET", "/api/servers", nil)
	req.Header.Set("X-API-Key", "secret")
	rr = httptest.NewRecorder()
	srv.Handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 with valid API key, got %d", rr.Code)
	}
}
