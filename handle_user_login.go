package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Lanrey-waju/gChirpy/internal/users"
)

func (cfg *apiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "wrong http method", http.StatusMethodNotAllowed)
	}

	type payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	pd := payload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pd)
	if err != nil {
		log.Printf("error decoding JSON: %w", err)
		return
	}
	user, err := cfg.DB.CheckUser(pd.Email, pd.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist or wrong credentials")
		return
	}

	userInfo := users.ReturnUserVal{
		ID:    user.ID,
		Email: user.Email,
	}
	log.Printf("Logged in Successfully!")
	respondWithJSON(w, http.StatusOK, userInfo)

}
