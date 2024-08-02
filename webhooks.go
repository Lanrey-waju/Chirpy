package main

import (
	"encoding/json"
	"net/http"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string         `json:"event"`
		Data  map[string]int `jaon:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't retrieve API key")
		return
	}

	if apiKey != cfg.apiKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.DB.UpgradeUser(params.Data["user_id"])
	if err != nil {
		respondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
