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
		return
	}

	type payload struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_Seconds *int   `json:"expires_in_seconds"`
	}

	pd := payload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pd)
	if err != nil {
		log.Printf("error decoding JSON: %v", err)
		respondWithError(w, http.StatusBadRequest, "bad request")
		return
	}
	if pd.Expires_in_Seconds == nil || *pd.Expires_in_Seconds > 86400 {
		pd.Expires_in_Seconds = new(int)
		*pd.Expires_in_Seconds = 86400
	}

	user, err := cfg.DB.GetUserByEmail(pd.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist or wrong credentials")
		return
	}

	token := CreateJWT(user.ID, *pd.Expires_in_Seconds)

	userInfo := users.ReturnUserVal{
		ID:    user.ID,
		Email: user.Email,
		Token: token,
	}
	log.Printf("Logged in Successfully!")
	respondWithJSON(w, http.StatusOK, userInfo)

}
