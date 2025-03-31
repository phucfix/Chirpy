package main

import (
	"net/http"
	"log"
	"sync/atomic"
	"fmt"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
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
	serveMux.HandleFunc("/healthz", handlerReadiness)

	// Add handler for number of requests
	serveMux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	// Add handler for reset number of hits
	serveMux.HandleFunc("/reset", apiCfg.handlerResetHit)

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}
