package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Lanrey-waju/gChirpy/internal/database"
)

type apiConfig struct {
	DB             *database.DB
	fileserverHits int
}

func main() {
	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}
	apiCfg := apiConfig{
		DB:             db,
		fileserverHits: 0,
	}

	handleRequests(&apiCfg)
}

func handleRequests(cfg *apiConfig) {
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/app/*", http.StripPrefix("/app", cfg.middlewareMetrics(http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("/admin/metrics", cfg.noOfRequests)
	mux.HandleFunc("/api/reset", cfg.reset)
	mux.HandleFunc("GET /api/healthz", ready)
	mux.HandleFunc("/api/chirps", cfg.ChirpsHandler)
	mux.HandleFunc("/api/chirps/{id}", cfg.GetSingleChirpHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files on %s:%s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) noOfRequests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(200)
	hits := cfg.fileserverHits
	fmt.Fprintf(
		w,
		`<html>

<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>

</html>`,
		hits,
	)
}
