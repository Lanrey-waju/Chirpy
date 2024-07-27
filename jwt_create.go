package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func CreateJWT(id int, expiration_time int) string {

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading environment variables: %v", err)
		return ""
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))
	if len(signingKey) == 0 {
		log.Println("secret key not set")
		return ""
	}

	// create claims
	claims := &jwt.RegisteredClaims{

		Issuer:    "chirpy",
		Subject:   strconv.Itoa(id),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(expiration_time))),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("error signing token: %v", err)
		return ""
	}
	return ss
}
