package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HealthRoute() chi.Router {
	router := chi.NewRouter()
	router.Method(http.MethodPost, "/all", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode("Hello world")
	}))
	return router
}
