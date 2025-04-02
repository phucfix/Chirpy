package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func respondWithError(w http.ResponseWriter, code int,	msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}
	
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
