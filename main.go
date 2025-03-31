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
	serveMux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	httpServer := &http.Server{
		Handler: serveMux,
		Addr: ":" + port,
	}


	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(httpServer.ListenAndServe())
}
