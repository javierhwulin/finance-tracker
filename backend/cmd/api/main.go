package main

import (
	"log/slog"
	"net/http"
	"os"

	httpRouter "github.com/javierhwulin/finance-tracker/internal/http"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	router := httpRouter.NewRouter()

	logger.Info("Starting server", "port", ":8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Error("Error init server", "error", err)
	}
}
