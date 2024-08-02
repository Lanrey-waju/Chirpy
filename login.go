package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
	"github.com/Lanrey-waju/gChirpy/internal/users"
)

func (cfg *apiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		users.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	if r.Method != http.MethodPost {
		http.Error(w, "wrong http method", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("error decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist or wrong credentials")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	// Generate refresh tokens
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error creating refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	err = cfg.DB.SaveRefreshToken(user.ID, refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: users.User{
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
