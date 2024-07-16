package main

import (
	"fmt"
	"log"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}
	mux.Handle("/app/*", http.StripPrefix("/app", apiCfg.middlewareMetrics(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /metrics", apiCfg.noOfRequests)
	mux.HandleFunc("/reset", apiCfg.reset)
	mux.HandleFunc("GET /healthz", ready)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files on %s:%s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (c *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (c *apiConfig) noOfRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(200)
	hits := c.fileserverHits
	fmt.Fprintf(w, "Hits: %d", hits)
}
