package handlers

import (
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
)

type Router struct {
	cfg *config.ApiConfig
}

func RegisterHandlers(cfg *config.ApiConfig, mux *http.ServeMux) {
	router := &Router{cfg: cfg}
	mux.HandleFunc("/healthz", router.HealthCheck)
}

func (router *Router) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
