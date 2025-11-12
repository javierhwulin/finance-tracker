package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/javierhwulin/finance-tracker/internal/app"
	"github.com/javierhwulin/finance-tracker/internal/config"
	httpRouter "github.com/javierhwulin/finance-tracker/internal/http"
	"github.com/javierhwulin/finance-tracker/internal/http/dto"
	"github.com/javierhwulin/finance-tracker/internal/repo"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create test config
	cfg := &config.Config{
		Version: "1.0.0",
		Port:    "8080",
		Env:     "test",
	}

	// Setup test dependencies
	userRepo := repo.NewUserMemoryRepository()
	application := app.NewApp(userRepo)
	router := httpRouter.NewRouter(cfg, application)

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

func TestCreateUser(t *testing.T) {
	// Setup
	cfg := &config.Config{Version: "1.0.0", Port: "8080", Env: "test"}
	userRepo := repo.NewUserMemoryRepository()
	application := app.NewApp(userRepo)
	router := httpRouter.NewRouter(cfg, application)

	// Test valid user creation
	t.Run("valid user", func(t *testing.T) {
		payload := dto.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		var response dto.CreateUserResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatal(err)
		}

		if response.Email != payload.Email {
			t.Errorf("handler returned unexpected email: got %v want %v", response.Email, payload.Email)
		}
	})

	// Test duplicate email
	t.Run("duplicate email", func(t *testing.T) {
		payload := dto.CreateUserRequest{
			Email:    "test@example.com",
			Password: "password456",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	// Test invalid email
	t.Run("invalid email", func(t *testing.T) {
		payload := dto.CreateUserRequest{
			Email:    "notanemail",
			Password: "password123",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})
}
