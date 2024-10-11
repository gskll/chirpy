package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/gskll/chirpy2/internal/chirp"
	"github.com/gskll/chirpy2/internal/config"
	"github.com/gskll/chirpy2/internal/database"
	"github.com/gskll/chirpy2/internal/user"
)

type APIRouter struct {
	cfg *config.ApiConfig
}

func RegisterAPIHandlers(prefix string, cfg *config.ApiConfig, mux *http.ServeMux) {
	router := APIRouter{cfg: cfg}

	mux.HandleFunc("GET "+prefix+"/healthz", router.HealthCheck)
	mux.HandleFunc("POST "+prefix+"/users", router.CreateUser)
	mux.HandleFunc("POST "+prefix+"/chirps", router.CreateChirp)
	mux.HandleFunc("GET "+prefix+"/chirps", router.GetChirps)
	mux.HandleFunc("GET "+prefix+"/chirps/{chirpID}", router.GetChirp)
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

func (router *APIRouter) CreateChirp(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
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

	dbChirp, err := router.cfg.Db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{Body: cleaned, UserID: params.UserId},
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
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp := chirp.NewChirp(dbChirp)

	respondWithJSON(w, http.StatusOK, chirp)
}
