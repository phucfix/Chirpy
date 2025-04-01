package main

import (
	"strings"
	"encoding/json"
	"net/http"
)

type Chirp struct {
	Body string `json:"body"`
}

func handleValidateChirp(w http.ResponseWriter, req *http.Request) {
	type Response struct {
		CleanedBody string `json:"cleaned_body"`
		Valid 		bool   `json:"valid"`
	}

	// Decode json req body
	decoder := json.NewDecoder(req.Body)
	chirpReq := Chirp{}
	err := decoder.Decode(&chirpReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding json")
		return
	}


	// Validation
	if len(chirpReq.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := replaceBadWords(chirpReq.Body)
	respondWithJSON(w, http.StatusOK, Response{ CleanedBody: cleaned, Valid: true })
}

func replaceBadWords(sentence string) string {
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(sentence, " ")
	for i := range words {
		for j := range badWords {
			if strings.ToLower(words[i]) == badWords[j] {
				words[i] = "****"
				break
			}
		}
	}

	return strings.Join(words, " ")
}
