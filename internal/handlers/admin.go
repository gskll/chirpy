package handlers

import (
	"fmt"
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
)

type AdminRouter struct {
	cfg *config.ApiConfig
}

func RegisterAdminHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := &AdminRouter{cfg: cfg}
	mux.HandleFunc("GET "+prefix+"/metrics", router.GetMetrics)
	mux.HandleFunc("POST "+prefix+"/reset", router.ResetMetrics)
}

func (router *AdminRouter) GetMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	body := fmt.Sprintf(
		"<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		router.cfg.FileServerHits.Load(),
	)
	w.Write([]byte(body))
}

func (router *AdminRouter) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	if router.cfg.Platform != config.DEV {
		respondWithError(w, http.StatusForbidden, "reset forbidden on current platform")
		return
	}

	router.cfg.Db.DeleteUsers(r.Context())
	router.cfg.FileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("App reset"))
}
