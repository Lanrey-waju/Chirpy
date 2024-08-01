package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/Lanrey-waju/gChirpy/internal/users"
	"golang.org/x/crypto/bcrypt"
)

var ErrAlreadyExists = errors.New("already exists")

func (db *DB) GetUsers() ([]users.User, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbData := DBStructure{
		Users: make(map[int]users.User),
	}
	var users []users.User

	dat, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			return users, nil
		}
		return users, err
	}
	if err := json.Unmarshal(dat, &dbData); err != nil {
		return users, err
	}

	for _, user := range dbData.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) CreateUser(email string, hashedPassword string) (users.User, error) {
	if _, err := db.GetUserByEmail(email); errors.Is(err, ErrNotExist) {
		return users.User{}, ErrAlreadyExists
	}

	dbStructure, err := db.loadDB()
	if err != nil {
		return users.User{}, err
	}

	id := len(dbStructure.Users) + 1
	user := users.User{
		ID:             id,
		Email:          email,
		HashedPassword: hashedPassword,
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return users.User{}, err
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (users.User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return users.User{}, nil
	}
	for _, user := range dbData.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return users.User{}, errors.New("user does not exist")
}

func (db *DB) GetUserByID(id int) (users.User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return users.User{}, err
	}
	if user, ok := dbData.Users[id]; ok {
		return user, nil
	}
	return users.User{}, errors.New("user does not exist")
}

func (db *DB) UpdateUser(id int, email, hashedPaassword string) (users.User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return users.User{}, nil
	}
	user, err := db.GetUserByID(id)
	if err != nil {
		return users.User{}, err
	}

	user.Email = email
	user.HashedPassword = hashedPaassword

	dbData.Users[user.ID] = user

	if err := db.writeDB(dbData); err != nil {
		log.Printf("Error writing to database: %v", err)
		return users.User{}, err
	}

	return user, nil

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
