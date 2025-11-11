package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

const version = "1.0.0"

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": version})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	http.HandleFunc("/api/health", healthCheckHandler)

	logger.Info("Starting server", "port", ":8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("Error init server", "error", err)
	}
}
