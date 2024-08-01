package main

import (
	"log"
	"net/http"
	"sort"
)

func (cfg *apiConfig) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	log.Println("Got here!")
	users, err := cfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	respondWithJSON(w, http.StatusOK, users)

}
