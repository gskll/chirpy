package handlers

import (
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
)

type APIRouter struct {
	cfg *config.ApiConfig
}

func RegisterAPIHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := APIRouter{cfg: cfg}

	mux.HandleFunc("GET "+prefix+"/healthz", router.HealthCheck)
}

func (router *APIRouter) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
