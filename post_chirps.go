package main

import (
	"encoding/json"
	"net/http"
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

	chirp, err := cfg.DB.CreateChirp(cleanChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}
