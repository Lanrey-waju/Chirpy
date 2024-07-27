package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/Lanrey-waju/gChirpy/internal/users"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp      `json:"chirps"`
	Users  map[int]users.User `json:"users"`
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  new(sync.RWMutex),
	}

	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil

}

func (db *DB) ensureDB() error {
	initialData := `{"chirps": {}, "users": {}}`

	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		err := os.WriteFile(db.path, []byte(initialData), 0644)
		if err != nil {
			return err
		}
	}
	return nil

}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbData.Chirps) + 1
	chirp := Chirp{
		ID:   newID,
		Body: body,
	}

	dbData.Chirps[newID] = chirp

	if err := db.writeDB(dbData); err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	var dbData DBStructure
	var chirps []Chirp

	dat, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Can't read file")
		}
		return chirps, err
	}
	if err := json.Unmarshal(dat, &dbData); err != nil {
		return chirps, err
	}

	for _, chirp := range dbData.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	readData, err := os.ReadFile(db.path)
	dbData := DBStructure{
		Chirps: make(map[int]Chirp),
		Users:  make(map[int]users.User),
	}
	if err != nil {
		if os.IsNotExist(err) {
			return dbData, nil
		}
		return dbData, err
	}
	if err := json.Unmarshal(readData, &dbData); err != nil {
		return dbData, err
	}
	return dbData, nil

}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	if err := os.WriteFile(db.path, data, 0644); err != nil {
		return err
	}
	return nil
}

func (db *DB) GetSingleChirp(id int) (Chirp, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if chirp, ok := dbData.Chirps[id]; ok {
		return chirp, nil
	}
	return Chirp{}, errors.New("chirp not found")
}

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

func (db *DB) CreateUser(email string, password string) (users.User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		log.Printf("Error loading database: %v", err)
		return users.User{}, err
	}

	newID := len(dbData.Users) + 1
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return users.User{}, err
	}
	user := users.User{
		ID:       newID,
		Email:    email,
		Password: hashedPassword,
	}

	dbData.Users[newID] = user

	if err := db.writeDB(dbData); err != nil {
		log.Printf("Error writing to database: %v", err)
		return users.User{}, err
	}
	log.Println("Database successfully written")

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
		return users.User{}, nil
	}
	if user, ok := dbData.Users[id]; ok {
		return user, nil
	}
	return users.User{}, errors.New("user does not exist")
}

func (db *DB) UpdateUser(id int, email, password string) (users.User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		return users.User{}, nil
	}
	user, err := db.GetUserByID(id)
	if err != nil {
		return users.User{}, err
	}

	user.Email = email
	user.Password, err = hashPassword(password)
	if err != nil {
		return users.User{}, err // Handle the hashing error
	}

	updatedUser := users.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}

	dbData.Users[user.ID] = updatedUser

	if err := db.writeDB(dbData); err != nil {
		log.Printf("Error writing to database: %v", err)
		return users.User{}, err
	}

	return updatedUser, nil

}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
