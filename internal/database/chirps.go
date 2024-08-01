package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
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

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	log.Println("Debugging: Here after")
	if err != nil {
		return Chirp{}, err
	}

	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   newID,
		Body: body,
	}

	dbStructure.Chirps[newID] = chirp

	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}
	return chirp, nil
}
