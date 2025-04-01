package main

import (
	"encoding/json"
	"log"
	"net/http"
)



func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type ChirpRequest struct {
		Body string `json:"body"`
	}

	w.Header().Set("Content-Type", "application/json")

	// Decode json req body
	decoder := json.NewDecoder(req.Body)
	chirpReq := ChirpRequest{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		log.Printf("Error decoding Chirp's requests: %v", err)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Something went wrong"}`))
		return
	}

	// Validation
	if len(chirpReq.Body) > 140 {
		log.Printf(`Sending response: {"error":"Chirp is too long"}`)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"Chirp is too long"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"valid":true}`))
}

