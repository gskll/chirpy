package handlers

import (
	"encoding/json"
	"fmt"
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
	w.Header().Set("Content-Type", "application/json")

	type reqParams struct {
		Body string `json:"body"`
	}

	resp := make(map[string]any)

	decoder := json.NewDecoder(r.Body)
	params := reqParams{}
	err := decoder.Decode(&params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp["error"] = fmt.Sprintf("Error decoding params: %s", err)
		dat, err := json.Marshal(resp)
		if err != nil {
			return
		}
		w.Write(dat)
		return
	}

	if len(params.Body) > 140 {
		resp["error"] = fmt.Sprintf("Chirp is too long. Max chars 140. Actual: %d", len(params.Body))
		dat, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(dat)
		return
	}

	resp["valid"] = true
	dat, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
	return
}
