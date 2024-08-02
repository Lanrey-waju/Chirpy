package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Lanrey-waju/gChirpy/internal/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB             *database.DB
	fileserverHits int
	jwtSecret      string
	apiKey         string
}

func main() {

	godotenv.Load(".env")

	jwtSecret := os.Getenv("JWT_SECRET")
	apiKey := os.Getenv("API_KEY")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	db, err := database.NewDB("database.json")
	if err != nil {
		log.Fatalf("Error creating database: %v", err)
	}

	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if dbg != nil && *dbg {
		err := db.ResetDB()
		if err != nil {
			log.Fatal(err)
		}
	}
	apiCfg := apiConfig{
		DB:             db,
		fileserverHits: 0,
		jwtSecret:      jwtSecret,
		apiKey:         apiKey,
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
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirp)

	mux.HandleFunc("/api/users", cfg.UsersHandler)

	mux.HandleFunc("POST /api/refresh", cfg.HandleTokenRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.HandleTokenRevoke)

	mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerWebhook)

	mux.HandleFunc("/api/login", cfg.LoginUser)
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
