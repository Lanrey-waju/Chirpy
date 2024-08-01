package database

import (
	"log"
	"time"

	"github.com/Lanrey-waju/gChirpy/internal/users"
)

type RefreshToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"body"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (db *DB) SaveRefreshToken(id int, token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil
	}

	refreshToken := RefreshToken{
		UserID:    id,
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	dbStructure.RefreshTokens[token] = refreshToken

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}
	return nil

}

func (db *DB) CheckRefreshToken(token string) (users.User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return users.User{}, err
	}

	refreshToken, ok := dbStructure.RefreshTokens[token]
	if !ok {
		return users.User{}, ErrNotExist
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return users.User{}, ErrNotExist
	}

	user, err := db.GetUserByID(refreshToken.UserID)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	delete(dbStructure.RefreshTokens, token)

	if err := db.writeDB(dbStructure); err != nil {
		log.Printf("Error writing to database: %v", err)
		return err
	}

	return nil

}
