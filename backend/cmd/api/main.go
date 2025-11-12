package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/javierhwulin/finance-tracker/internal/config"
	httpRouter "github.com/javierhwulin/finance-tracker/internal/http"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create router with dependencies
	router := httpRouter.NewRouter(cfg)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Starting server", "port", addr, "env", cfg.Env, "version", cfg.Version)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
