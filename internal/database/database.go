package database

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
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

func ChirpsHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	switch r.Method {
	case http.MethodGet:
		GetChirpHandler(w, r, db)
	case http.MethodPost:
		PostChirpHandler(w, r, db)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func PostChirpHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	type payload struct {
		Body string `json:"body"`
	}

	pd := payload{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pd)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp body")
		return
	}

	const chirpmaxlength = 140
	if len(pd.Body) == 0 || len(pd.Body) > chirpmaxlength {
		respondWithError(w, http.StatusBadRequest, "Chirp is invalid")
		return
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	cleanChirp := removeProfaneWords(profaneWords, pd.Body)

	chirp, err := db.CreateChirp(cleanChirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}

func GetChirpHandler(w http.ResponseWriter, r *http.Request, db *DB) {
	chirps, err := db.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps")
		return
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJSON(w, http.StatusOK, chirps)

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
