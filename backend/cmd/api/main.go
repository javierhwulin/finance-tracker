package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func greet(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World! %s", time.Now())
	if err != nil {
		return
	}
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	http.HandleFunc("/", greet)

	logger.Info("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("Error init server", err)
	}
}
