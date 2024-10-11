package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/gskll/chirpy2/internal/auth"
	"github.com/gskll/chirpy2/internal/chirp"
	"github.com/gskll/chirpy2/internal/config"
	"github.com/gskll/chirpy2/internal/database"
)

type APIRouter struct {
	cfg *config.ApiConfig
}

func RegisterAPIHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := APIRouter{cfg: cfg}

	mux.HandleFunc("GET "+prefix+"/healthz", router.HealthCheck)

	mux.HandleFunc("POST "+prefix+"/users", router.CreateUser)
	mux.HandleFunc("POST "+prefix+"/login", router.LoginUser)
	mux.HandleFunc("POST "+prefix+"/refresh", router.RefreshToken)
	mux.HandleFunc("POST "+prefix+"/revoke", router.RevokeRefreshToken)

	mux.HandleFunc("POST "+prefix+"/chirps", router.CreateChirp)
	mux.HandleFunc("GET "+prefix+"/chirps", router.GetChirps)
	mux.HandleFunc("GET "+prefix+"/chirps/{chirpID}", router.GetChirp)
}

func (router *APIRouter) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (router *APIRouter) CreateChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userId, err := auth.ValidateJWT(token, router.cfg.JWTSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type reqParams struct {
		Body string `json:"body"`
	}

	params := reqParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = chirp.ValidateLength(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	cleaned := chirp.Clean(params.Body)

	dbChirp, err := router.cfg.Db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{Body: cleaned, UserID: userId},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp := chirp.NewChirp(dbChirp)

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (router *APIRouter) GetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := router.cfg.Db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirps := make([]chirp.Chirp, 0, len(dbChirps))
	for _, dbChirp := range dbChirps {
		chirp := chirp.NewChirp(dbChirp)
		chirps = append(chirps, chirp)
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (router *APIRouter) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}
	dbChirp, err := router.cfg.Db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp := chirp.NewChirp(dbChirp)

	respondWithJSON(w, http.StatusOK, chirp)
}
