package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	type Chirp struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding chirp body")
		return
	}

	const chirpmaxlength = 140
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	if len(chirp.Body) > chirpmaxlength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, returnVal{
		CleanedBody: removeProfaneWords(profaneWords, chirp.Body),
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
