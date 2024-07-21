package main

import "net/http"

func (c *apiConfig) reset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	c.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File server hits reset to 0"))
}
