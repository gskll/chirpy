package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gskll/chirpy2/internal/chirp"
	"github.com/gskll/chirpy2/internal/config"
	"github.com/gskll/chirpy2/internal/user"
)

type APIRouter struct {
	cfg *config.ApiConfig
}

func RegisterAPIHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := APIRouter{cfg: cfg}

	mux.HandleFunc("GET "+prefix+"/healthz", router.HealthCheck)
	mux.HandleFunc("POST "+prefix+"/users", router.CreateUser)
	mux.HandleFunc("POST "+prefix+"/validate_chirp", router.ValidateChirpLength)
}

func (router *APIRouter) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (router *APIRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email string `json:"email"`
	}

	params := reqParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := router.cfg.Db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := user.NewUser(dbUser)
	respondWithJSON(w, http.StatusCreated, user)
}

func (router *APIRouter) ValidateChirpLength(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body string `json:"body"`
	}

	params := reqParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err := chirp.ValidateLength(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cleaned := chirp.Clean(params.Body)

	respondWithJSON(w, http.StatusOK, map[string]any{"cleaned_body": cleaned})
}
