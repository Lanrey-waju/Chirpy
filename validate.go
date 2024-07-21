package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rivo/uniseg"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	type Chirp struct {
		Body string `json:"body"`
	}

	type returnVal struct {
		Error string `json:"error"`
	}

	type Validity struct {
		Valid bool `json:"valid"`
	}
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respBody := returnVal{
			Error: "Something went wrong",
		}
		w.WriteHeader(500)
		dat, _ := json.Marshal(respBody)
		w.Write(dat)
		return
	}

	fmt.Println(uniseg.GraphemeClusterCount(chirp.Body))
	if uniseg.GraphemeClusterCount(chirp.Body) > 140 {
		respBody := returnVal{
			Error: "Chirpy is too long",
		}
		dat, _ := json.Marshal(respBody)
		w.WriteHeader(400)
		w.Write(dat)
		return
	}
	w.WriteHeader(200)
	respBody := Validity{
		Valid: true,
	}
	dat, _ := json.Marshal(respBody)
	w.Write(dat)

}
