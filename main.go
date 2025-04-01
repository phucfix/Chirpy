package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"os"
	"database/sql"
	_ "github.com/lib/pq"

	"github.com/phucfix/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries	   *database.Queries
}

func main() {
	dbURL := os.Getenv("DB_URL")
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
	}

	serveMux := http.NewServeMux()
	// Add handler for file server
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	// Add handler for the readiness
	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)

	// Add handler for number of requests
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// Add handler for reset number of hits
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetHit)
	serveMux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}
