package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gskll/chirpy2/internal/config"
)

type APIRouter struct {
	cfg *config.ApiConfig
}

func RegisterAPIHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := APIRouter{cfg: cfg}

	mux.HandleFunc("GET "+prefix+"/healthz", router.HealthCheck)
	mux.HandleFunc("POST "+prefix+"/validate_chirp", router.ValidateChirpLength)
}

func (router *APIRouter) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (router *APIRouter) ValidateChirpLength(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
	}

	params := reqParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long. Max 140 chars")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]any{"valid": true})
}
