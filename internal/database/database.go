package database

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
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

func removeProfaneWords(profaneWords []string, body string) string {
	substrings := strings.Split(body, " ")
	for i, substring := range substrings {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(substring) == strings.ToLower(profaneWord) {
				substrings[i] = "****"
			}
		}
	}
	return strings.Join(substrings, " ")
}

func (db *DB) ensureDB() error {
	initialData := `{"chirps": {}}`

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
