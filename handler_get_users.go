package main

import (
	"net/http"
	"sort"
)

func (cfg *apiConfig) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
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
