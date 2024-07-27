package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func (cfg *apiConfig) HandlePutUser(w http.ResponseWriter, r *http.Request) {
	pd := payload{}

	if err := json.NewDecoder(r.Body).Decode(&pd); err != nil {
		return
	}

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading environment variables: %v", err)
		return
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))
	if len(signingKey) == 0 {
		log.Println("secret key not set")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization Header is Missing")
		return
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "invalid or expired token", http.StatusUnauthorized)
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		http.Error(w, "invalid authorization header", http.StatusUnauthorized)
	}

	cfg.DB.UpdateUser(userID, pd.Email, pd.Password)

	resp := map[string]interface{}{
		"id":    userID,
		"email": pd.Email,
	}
	respondWithJSON(w, http.StatusOK, resp)

}
