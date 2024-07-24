package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) GetSingleChirpHandler(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("id")
	id, err := strconv.Atoi(value)
	if err != nil {
		http.Error(w, "error converting string to int", http.StatusInternalServerError)
	}
	chirp, err := cfg.DB.GetSingleChirp(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)

}
