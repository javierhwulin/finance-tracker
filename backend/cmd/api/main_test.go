package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/javierhwulin/finance-tracker/internal/config"
	httpRouter "github.com/javierhwulin/finance-tracker/internal/http"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create test config directly without parsing flags
	cfg := &config.Config{
		Version: "1.0.0",
		Port:    "8080",
		Env:     "test",
	}

	router := httpRouter.NewRouter(cfg)

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]string
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	if response["status"] != "ok" {
		t.Errorf("handler returned unexpected status: got %v want %v", response["status"], "ok")
	}

	if response["version"] != cfg.Version {
		t.Errorf("handler returned unexpected version: got %v want %v", response["version"], cfg.Version)
	}
}
