package main

import "net/http"

func (c *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	c.fileserverHits = 0
	w.Write([]byte("File server hits reset to 0"))
}
