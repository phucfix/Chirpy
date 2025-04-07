package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"

	"github.com/phucfix/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries	   *database.Queries
	platform       string
	jwtSecret	   string
}

func main() {
	// Load enviroment variable
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	
	platform := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	// Open a connection to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't open connection to the database: %w", err)
	}
	// Create a new *database.Queries
	dbQueries := database.New(db)

	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
	}

	serveMux := http.NewServeMux()
	// Add handler for file server
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Add handler for the readiness
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Add handler for number of requests
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// Add handler for reset number of hits
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerDeleteUser)

	// User API
	serveMux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	// Chirps API
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serveMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)

	// Token
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}
