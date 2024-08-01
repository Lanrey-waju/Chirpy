package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/Lanrey-waju/gChirpy/internal/users"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]users.User      `json:"users"`
	RefreshTokens map[string]RefreshToken `json:"refresh_tokens"`
}

var ErrNotExist = errors.New("resosurce does not exist")

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  new(sync.RWMutex),
	}

	err := db.ensureDB()

	return db, err

}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps:        map[int]Chirp{},
		Users:         map[int]users.User{},
		RefreshTokens: map[string]RefreshToken{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}
	if err := os.WriteFile(db.path, data, 0600); err != nil {
		return err
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}
	readData, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	if err := json.Unmarshal(readData, &dbStructure); err != nil {
		return dbStructure, err
	}
	return dbStructure, nil

}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) ResetDB() error {
	err := os.Remove(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return db.ensureDB()
}
