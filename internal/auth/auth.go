package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var ErrNoAuthHeaderIncluded = errors.New("Header not included")

func MakeJWT(id int) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Printf("error loading environment variables: %v", err)
		return "", err
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))
	if len(signingKey) == 0 {
		log.Println("secret key not set")
		return "", err
	}

	// create claims
	claims := &jwt.RegisteredClaims{

		Issuer:    "chirpy",
		Subject:   strconv.Itoa(id),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(signingKey)
	if err != nil {
		log.Printf("error signing token: %v", err)
		return "", err
	}
	return ss, nil
}

// Generates a refresh toke struct of body amd issue time
func CreateRefreshToken() (string, error) {

	const numBytes = 32

	randomBytes := make([]byte, numBytes)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Println("error:", err)
		return "", err
	}
	refreshToken := hex.EncodeToString(randomBytes)

	return refreshToken, nil
}

func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "Bearer" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// HashPassword -
func HashPassword(password string) (string, error) {
	dat, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

func ValidateJWT(tokenString, tokenSecret string) (string, error) {
	claimsStruct := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claimsStruct,
		func(token *jwt.Token) (interface{}, error) { return []byte(tokenSecret), nil },
	)
	if err != nil {
		return "", err
	}

	userIDString, err := token.Claims.GetSubject()
	if err != nil {
		return "", err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		return "", err
	}

	if issuer != string("chirpy") {
		return "", errors.New("invalid issuer")
	}

	return userIDString, nil
}
