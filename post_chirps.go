package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
)

func (cfg *apiConfig) PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	type payload struct {
		Body string `json:"body"`
	}

	pd := payload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp body")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token")
		return
	}

	userIDString, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't validate token")
		return
	}

	userIDInt, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "server error")
		return
	}

	const chirpmaxlength = 140
	if len(pd.Body) == 0 || len(pd.Body) > chirpmaxlength {
		respondWithError(w, http.StatusBadRequest, "Chirp is invalid")
		return
	}

	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleanChirp := removeProfaneWords(profaneWords, pd.Body)

	chirp, err := cfg.DB.CreateChirp(userIDInt, cleanChirp)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}
