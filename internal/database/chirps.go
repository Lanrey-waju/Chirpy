package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	var dbStructure DBStructure
	var chirps []Chirp

	dat, err := os.ReadFile(db.path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalln("Can't read file")
		}
		return chirps, err
	}
	if err := json.Unmarshal(dat, &dbStructure); err != nil {
		return chirps, err
	}

	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetSingleChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	if chirp, ok := dbStructure.Chirps[id]; ok {
		return chirp, nil
	}
	return Chirp{}, errors.New("chirp not found")
}

func (db *DB) CreateChirp(authorID int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	user, err := db.GetUserByID(authorID)
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       newID,
		Body:     body,
		AuthorID: user.ID,
	}

	dbStructure.Chirps[newID] = chirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}
