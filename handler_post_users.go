package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Lanrey-waju/gChirpy/internal/users"
)

func (cfg *apiConfig) PostUsersHandler(w http.ResponseWriter, r *http.Request) {

	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	pd := payload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pd)
	if err != nil {
		log.Println("Error decoding JSON")
		respondWithError(w, http.StatusInternalServerError, "Error decoding user email")
		return
	}

	user, err := cfg.DB.CreateUser(pd.Email, pd.Password)
	if err != nil {
		log.Println("Error creating user")
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	log.Printf("Created user: %+v", user)

	userInfo := users.ReturnUserVal{
		ID:    user.ID,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusCreated, userInfo)
	log.Println("Response sent")
}
