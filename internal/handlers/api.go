package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gskll/chirpy2/internal/auth"
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

func (router *APIRouter) RefreshToken(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	dbToken, err := router.cfg.Db.GetRefreshToken(r.Context(), rToken)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if dbToken.ExpiresAt.Before(time.Now()) || dbToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "token expired")
		return
	}

	token, err := auth.MakeJWT(dbToken.UserID, router.cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (router *APIRouter) RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	err = router.cfg.Db.RevokeRefreshToken(r.Context(), rToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (router *APIRouter) CreateUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := reqParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbUser, err := router.cfg.Db.CreateUser(
		r.Context(),
		database.CreateUserParams{HashedPassword: hashedPassword, Email: params.Email},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := user.NewUser(dbUser)
	respondWithJSON(w, http.StatusCreated, user)
}

func (router *APIRouter) LoginUser(w http.ResponseWriter, r *http.Request) {
	type reqParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := reqParams{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	dbUser, err := router.cfg.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, router.cfg.JWTSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = router.cfg.Db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{Token: refreshToken, UserID: dbUser.ID},
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user := user.NewUserWithTokens(dbUser, token, refreshToken)

	respondWithJSON(w, http.StatusOK, user)
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
