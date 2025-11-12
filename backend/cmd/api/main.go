package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/javierhwulin/finance-tracker/internal/app"
	"github.com/javierhwulin/finance-tracker/internal/config"
	httpRouter "github.com/javierhwulin/finance-tracker/internal/http"
	"github.com/javierhwulin/finance-tracker/internal/repo"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Initialize logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Initialize repositories
	userRepo := repo.NewUserMemoryRepository()

	// Wire up application
	application := app.NewApp(userRepo)

	// Create router with dependencies
	router := httpRouter.NewRouter(cfg, application)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("Starting server", "port", addr, "env", cfg.Env, "version", cfg.Version)

	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Error("Error starting server", "error", err)
		os.Exit(1)
	}
}
