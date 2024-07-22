package database

import (
	"encoding/json"
	"log"
	"net/http"
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

func PostChirpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	type returnVal struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp body")
		return
	}

	const chirpmaxlength = 140
	if len(chirp.Body) > chirpmaxlength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	respondWithJSON(w, http.StatusOK, returnVal{
		Body: removeProfaneWords(profaneWords, chirp.Body),
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("5XX error: %s", msg)
	}
	type errorRessponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorRessponse{
		Error: msg,
	})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write((dat))
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

	data, err := json.Marshal(dbData)
	if err != nil {
		return Chirp{}, err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		return Chirp{}, err
	}
	return chirp, nil
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
