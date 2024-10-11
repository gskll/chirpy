package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/gskll/chirpy2/internal/auth"
)

const UserUpgradedEvent = "user.upgraded"

func (router *APIRouter) UpgradeUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if apiKey != router.cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	params := struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Event != UserUpgradedEvent {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = router.cfg.Db.UpgradeUser(r.Context(), params.Data.UserId)
	if err != nil {
		handleDatabaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
