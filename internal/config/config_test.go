package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	content := `
api:
  port: 9090
  api_key: "test-key"
servers:
  - name: "server1"
    host: "192.168.1.100"
    username: "admin"
    password: "pass"
  - name: "server2"
    host: "192.168.1.101"
    port: 624
    username: "admin"
    password: "pass"
`

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.API.Port != 9090 {
		t.Errorf("expected port 9090, got %d", cfg.API.Port)
	}
	if cfg.API.APIKey != "test-key" {
		t.Errorf("expected api key 'test-key', got %q", cfg.API.APIKey)
	}
	if len(cfg.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(cfg.Servers))
	}
	if cfg.Servers[0].Port != 623 {
		t.Errorf("expected default port 623 for server1, got %d", cfg.Servers[0].Port)
	}
	if cfg.Servers[1].Port != 624 {
		t.Errorf("expected port 624 for server2, got %d", cfg.Servers[1].Port)
	}
}

func TestLoadDefaults(t *testing.T) {
	content := `
servers:
  - name: "server1"
    host: "10.0.0.1"
    username: "admin"
    password: "pass"
`

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.API.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.API.Port)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	if err := os.WriteFile(path, []byte("{{invalid yaml"), 0600); err != nil {
		t.Fatal(err)
	}

	_, err := Load(path)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestFindServer(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "server1", Host: "10.0.0.1"},
			{Name: "server2", Host: "10.0.0.2"},
		},
	}

	s, err := cfg.FindServer("server1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Host != "10.0.0.1" {
		t.Errorf("expected host 10.0.0.1, got %s", s.Host)
	}
}

func TestFindServerNotFound(t *testing.T) {
	cfg := &Config{
		Servers: []Server{
			{Name: "server1", Host: "10.0.0.1"},
		},
	}

	_, err := cfg.FindServer("nonexistent")
	if err == nil {
		t.Error("expected error for nonexistent server")
	}
}
