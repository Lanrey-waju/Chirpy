package main

import (
	"net/http"
	"strconv"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpIDInt, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirp ID")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode header")
		return
	}

	userODString, err := auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error validating token")
		return
	}

	userIDInt, err := strconv.Atoi(userODString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving user ")
		return
	}

	authorID, err := cfg.DB.GetChirpAuthorID(chirpIDInt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving chirp's author")
	}

	if userIDInt != authorID {
		respondWithError(w, http.StatusForbidden, "user is not chirps's author")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting userID")
		return
	}

	err = cfg.DB.DeleteChirp(chirpIDInt, userIDInt)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "error deleting chirp")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
