package main

import (
	"net/http"
	"log"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	// Add handler for the root path (/)
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	serveMux.HandleFunc("/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(http.StatusText(http.StatusOK)))
	})

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}

