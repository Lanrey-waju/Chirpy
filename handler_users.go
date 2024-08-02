package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Lanrey-waju/gChirpy/internal/auth"
)

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UsersHandler invoked with method: %s", r.Method) // Log request method

	switch r.Method {
	case http.MethodGet:
		cfg.HandleGetUser(w, r)
	case http.MethodPost:
		cfg.HandlePostUser(w, r)
	case http.MethodPut:
		cfg.handlerUsersUpdate(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (cfg *apiConfig) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Body != http.NoBody {
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err)
		return
	}

	user, err := cfg.DB.UserForRefreshToken(tokenString)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Println(err)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: accessToken,
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) HandleTokenRevoke(w http.ResponseWriter, r *http.Request) {
	if r.Body != http.NoBody {
		log.Println("Revoke request should not have a body")
		respondWithError(w, http.StatusBadRequest, "There should be no body in the request")
		return
	}

	tokenString, err := auth.GetBearerToken((r.Header))
	if err != nil {
		return
	}

	err = cfg.DB.RevokeToken(tokenString)
	if err != nil {
		log.Println("Error revoking token:", err)
		respondWithError(w, http.StatusUnauthorized, "not authorized")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
