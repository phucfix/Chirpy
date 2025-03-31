package main

import (
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
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

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}
