package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/javierhwulin/finance-tracker/internal/app"
	"github.com/javierhwulin/finance-tracker/internal/config"
	"github.com/javierhwulin/finance-tracker/internal/http/middleware"
)

// Reusable validator instance
var validate = validator.New()

// NewRouter creates and configures the HTTP router with all routes
func NewRouter(cfg *config.Config, app *app.App, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /api/health", healthCheckHandler(cfg))

	// User endpoints
	mux.HandleFunc("GET /api/users/:id", GetUserHandler(app.UserRepository))
	mux.HandleFunc("POST /api/users", CreateUserHandler(app.UserRepository))
	mux.HandleFunc("PUT /api/users/:id", UpdateUserHandler(app.UserRepository))
	mux.HandleFunc("DELETE /api/users/:id", DeleteUserHandler(app.UserRepository))

	// Middleware
	handler := middleware.Recovery(logger)(mux)
	handler = middleware.Logging(logger)(handler)
	handler = middleware.Cors()(handler)
	return handler
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
