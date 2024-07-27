package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) UsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UsersHandler invoked with method: %s", r.Method) // Log request method

	switch r.Method {
	case http.MethodGet:
		cfg.HandleGetUser(w, r)
	case http.MethodPost:
		cfg.HandlePostUser(w, r)
	case http.MethodPut:
		cfg.HandlePutUser(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
