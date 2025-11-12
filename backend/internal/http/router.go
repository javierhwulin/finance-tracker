package http

import (
	"encoding/json"
	"net/http"

	"github.com/javierhwulin/finance-tracker/internal/config"
)

// NewRouter creates and configures the HTTP router with all routes
func NewRouter(cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /api/health", healthCheckHandler(cfg))

	return mux
}

// healthCheckHandler returns the health check handler with config injected
func healthCheckHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := struct {
			Status  string `json:"status"`
			Version string `json:"version"`
		}{
			Status:  "ok",
			Version: cfg.Version,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}
