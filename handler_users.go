package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UsersHandler invoked with method: %s", r.Method) // Log request method

	switch r.Method {
	case http.MethodGet:
		cfg.GetUsersHandler(w, r)
	case http.MethodPost:
		cfg.PostUsersHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
