package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Query().Get("author_id") != "" {
		authorIDString := r.URL.Query().Get("author_id")
		authorIDInt, err := strconv.Atoi(authorIDString)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error getting author ID")
			return
		}
		chirps, err := cfg.DB.GetChirpsByAuthor(authorIDInt)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error retrieving author's chirps")
			return
		}
		switch r.URL.Query().Get("sort") {
		case "asc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].ID > chirps[j].ID
			})
		case "desc":
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].ID < chirps[j].ID
			})
		default:
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].ID < chirps[j].ID
			})

		}
		respondWithJSON(w, http.StatusOK, chirps)

	}
	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}

	switch r.URL.Query().Get("sort") {
	case "asc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	case "desc":
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	default:
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})

	}
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
