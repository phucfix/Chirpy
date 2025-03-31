package main

import (
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	httpServer := http.Server{
		Handler: serveMux,
		Addr: ":8080",
	}
	// Add handler for the root path (/)
	serveMux.Handle("/", http.FileServer(http.Dir(".")))

	httpServer.ListenAndServe()
}
