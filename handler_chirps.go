package main

import "net/http"

func (cfg *apiConfig) ChirpsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cfg.GetChirpHandler(w, r)
	case http.MethodPost:
		cfg.PostChirpHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
