package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
)

func (cfg *apiConfig) HandlePutUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Println(err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	if err != nil {
		log.Printf("error loading environment variables: %v", err)
		return
	}

	if len(cfg.jwtSecret) == 0 {
		log.Println("secret key not set")
		return
	}

	userIDString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, fmt.Sprint("Couldn't validate token: %w", err))
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	userID, err := strconv.Atoi(userIDString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse user ID")
		return
	}

	cfg.DB.UpdateUser(userID, params.Email, hashedPassword)

	resp := map[string]interface{}{
		"id":    userID,
		"email": params.Email,
	}
	respondWithJSON(w, http.StatusOK, resp)

}
