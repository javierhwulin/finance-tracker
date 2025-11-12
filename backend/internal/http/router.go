package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

const Version = "1.0.0"

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/health", healthCheckHandler).Methods("GET")
	return router
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"status": "ok", "version": Version})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
